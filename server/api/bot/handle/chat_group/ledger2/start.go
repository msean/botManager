package ledger2

import (
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type ResetLedgerHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	ShouldPermissionAwareWithOutAdmin
}

func (l *ResetLedgerHandler) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "上课" {
		return false
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *ResetLedgerHandler) Handle() error {

	today := time.Now().Format("2006-01-02")

	tx := global.GVA_MYSQL.Begin()

	// ============================
	// ✅ 1. 删除当天账单
	// ============================
	err := tx.Where(
		"bot_id = ? AND chat_group_id = ? AND work_date = ?",
		l.botModel.BotID,
		l.chatGroupID,
		today,
	).Delete(&ledger2.Ledger{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var session ledger2.LedgerSession

	err = tx.Where(
		"bot_id = ? AND chat_group_id = ? AND work_date = ?",
		l.botModel.BotID,
		l.chatGroupID,
		today,
	).First(&session).Error

	if err != nil {
		// 没有就创建
		session = ledger2.LedgerSession{
			BotID:       l.botModel.BotID,
			ChatGroupID: l.chatGroupID,
			WorkDate:    today,
			IsActive:    1,
		}

		if err := tx.Create(&session).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 有就重置
		session.IsActive = 1
		if err := tx.Save(&session).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return l.reply("🔄 今日账单已清空，已重新开始记账")
}

func (l *ResetLedgerHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}
	msg := botapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
