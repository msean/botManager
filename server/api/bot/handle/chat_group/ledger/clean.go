package ledger

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
)

type CleanHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	userID      string
	ShouldPermissionAwareWithOutAdmin
}

func (l *CleanHandler) Match(botModel bot.Bot, update tgbotapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "清空今日账单" {
		return false
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID
	return true
}

func (l *CleanHandler) Handle() error {

	now := time.Now()
	startOfDay := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	)
	endOfDay := startOfDay.Add(24 * time.Hour)

	result := global.GVA_MYSQL.
		Where("bot_id = ?", l.botModel.BotID).
		Where("chat_group_id = ?", l.chatGroupID).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Delete(&ledger.Ledger{})

	if result.Error != nil {
		return result.Error
	}

	return l.reply(fmt.Sprintf("✅ 已清空今日账单（共 %d 条）", result.RowsAffected))
}

func (l *CleanHandler) reply(text string) error {
	botSender, err := tgbotapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
