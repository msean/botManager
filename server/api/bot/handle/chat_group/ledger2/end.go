package ledger2

import (
	"fmt"
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type LedgerSummaryHandler struct {
	botModel    bot.Bot
	chatGroupID int64
}

func (l *LedgerSummaryHandler) Match(botModel bot.Bot, update botapi.Update) bool {

	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "下课" {
		return false
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *LedgerSummaryHandler) Handle() error {

	today := time.Now().Format("2006-01-02")

	var list []ledger2.Ledger

	// ✅ 查当天数据
	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ? AND DATE(created_at) = ?",
		l.botModel.BotID,
		l.chatGroupID,
		today,
	).Find(&list).Error

	if err != nil {
		return err
	}

	if len(list) == 0 {
		return l.reply("📊 今日暂无账单")
	}

	// ✅ 按 人+账号 分组
	type stat struct {
		in  float64
		out float64
	}

	groupMap := make(map[string]*stat)

	var totalIn float64
	var totalOut float64

	for _, v := range list {

		key := v.UserName + "+" + v.Account

		if _, ok := groupMap[key]; !ok {
			groupMap[key] = &stat{}
		}

		if v.ActionType == 1 {
			groupMap[key].in += v.Amount
			totalIn += v.Amount
		} else {
			groupMap[key].out += v.Amount
			totalOut += v.Amount
		}
	}

	// ✅ 拼接输出
	var builder strings.Builder

	builder.WriteString("📊 今日账单统计\n\n")

	for k, v := range groupMap {

		balance := v.in - v.out

		builder.WriteString(fmt.Sprintf(
			"%s\n收入：%.2f\n支出：%.2f\n余额：%.2f\n\n",
			k, v.in, v.out, balance,
		))
	}

	builder.WriteString("——————\n")
	builder.WriteString(fmt.Sprintf(
		"汇总：\n当日收入：%.2f\n当日支出：%.2f",
		totalIn, totalOut,
	))

	// ✅ 统计完顺便关闭（下课=结束）
	err = global.GVA_MYSQL.Model(&ledger2.LedgerSession{}).
		Where("bot_id = ? AND chat_group_id = ? AND work_date = ?",
			l.botModel.BotID, l.chatGroupID, today).
		Update("is_active", 0).Error

	if err != nil {
		return err
	}

	return l.reply(builder.String())
}

func (l *LedgerSummaryHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}
	msg := botapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
