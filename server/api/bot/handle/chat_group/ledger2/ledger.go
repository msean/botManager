package ledger2

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type LedgerRecordHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	msg         *botapi.Message

	userName string
	account  string
	amount   float64
	action   int
	rawText  string
}

func (l *LedgerRecordHandler) Match(botModel bot.Bot, update botapi.Update) bool {

	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	// ❗必须包含 + 和 金额
	if !strings.Contains(text, "+") {
		return false
	}

	// 👉 正则解析
	// 格式：章三+账号 +100 或 章三+账号-200 或 章三+余额+100
	re := regexp.MustCompile(`^(.+?)\+(.+?)\s*([+-])\s*(\d+(\.\d+)?)$`)
	match := re.FindStringSubmatch(text)

	if len(match) == 0 {
		return false
	}

	userName := strings.TrimSpace(match[1])
	account := strings.TrimSpace(match[2])
	sign := match[3]
	amountStr := match[4]

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return false
	}

	action := 1 // 默认收入

	// ✅ 核心：余额特殊处理
	if account == "余额" {
		action = 3
		l.account = ""
		if sign == "-" {
			amount = -amount
		}
	} else {
		// 正常收入/支出
		if sign == "-" {
			action = 2
		}
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID
	l.msg = update.Message

	l.userName = userName
	l.account = account
	l.amount = amount
	l.action = action
	l.rawText = text

	return true
}
func (l *LedgerRecordHandler) Handle() error {

	active, err := IsTodayActive(l.botModel.BotID, l.chatGroupID)
	if err != nil {
		return err
	}
	if !active {
		return nil
	}

	var count int64
	err = global.GVA_MYSQL.Model(&ledger2.Ledger{}).
		Where("message_id = ?", l.msg.MessageID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// ✅ 新增：账户唯一校验
	var exist ledger2.Ledger
	err = global.GVA_MYSQL.
		Where("bot_id = ? AND chat_group_id = ? AND user_name = ?", l.botModel.BotID, l.chatGroupID, l.userName).
		Order("id ASC").
		First(&exist).Error

	if err == nil {
		if exist.Account != l.account && l.account != "" {
			return l.reply(fmt.Sprintf(
				"❌ 账户错误：%s 已绑定账户【%s】，不能使用【%s】",
				l.userName,
				exist.Account,
				l.account,
			))
		}
	}

	record := ledger2.Ledger{
		OprUserID:        l.msg.From.ID,
		OprUserFirstName: l.msg.From.FirstName,
		OprUserLastName:  l.msg.From.LastName,
		OprUserNickname:  l.msg.From.UserName,

		ActionType: l.action,
		Amount:     l.amount,

		BotID:       l.botModel.BotID,
		UserName:    l.userName,
		Account:     l.account,
		ChatGroupID: l.chatGroupID,

		MessageID: l.msg.MessageID,
		RawInput:  l.rawText,
	}

	if err := global.GVA_MYSQL.Create(&record).Error; err != nil {
		return err
	}

	if l.action == 1 {
		return l.reply(fmt.Sprintf("✅ 收入记录成功：%s +%.2f", l.userName, l.amount))
	} else if l.action == 2 {
		return l.reply(fmt.Sprintf("✅ 支出记录成功：%s -%.2f", l.userName, l.amount))
	} else {
		return l.reply(fmt.Sprintf("✅ 余额调整成功：%s %+ .2f", l.userName, l.amount))
	}
}

func (l *LedgerRecordHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}
	msg := botapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
