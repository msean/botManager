package ledger

import (
	"errors"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
)

type LedgerShowFeeRateHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	ShouldPermissionAwareWithOutAdmin
}

func (l *LedgerShowFeeRateHandler) Match(botModel bot.Bot, update tgbotapi.Update) (match bool) {

	text := strings.TrimSpace(update.Message.Text)

	if text != "查看费率" {
		return false
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *LedgerShowFeeRateHandler) Handle() (err error) {
	if l.botModel.BotID == 0 || l.chatGroupID == 0 {
		return errors.New("bot or chat group invalid")
	}

	var permission ledger.LedgerPermission

	err = global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		l.botModel.BotID,
		l.chatGroupID,
	).First(&permission).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没设置过费率
			return l.reply("当前群聊尚未设置费率")
		}
		return err
	}

	return l.reply(
		"当前费率为：" + formatFeeRate(permission.CurrentFeeRate),
	)
}

func (l *LedgerShowFeeRateHandler) reply(text string) (err error) {
	var botSender *tgbotapi.BotAPI
	botSender, err = tgbotapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return
	}
	msg := tgbotapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}

func formatFeeRate(rate float64) string {
	// 如果你以后想加 %，可以在这里统一处理
	return strings.TrimRight(strings.TrimRight(
		strconv.FormatFloat(rate, 'f', 2, 64),
		"0",
	), ".")
}
