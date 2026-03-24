package ledger2

import (
	"errors"
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"gorm.io/gorm"
)

type StartLedgerHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	OnlyAdminAware
}

func (l *StartLedgerHandler) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "开始" {
		return false
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *StartLedgerHandler) Handle() error {
	var session ledger2.LedgerSession

	today := time.Now().Format("2006-01-02")

	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ? AND work_date = ?",
		l.botModel.BotID,
		l.chatGroupID,
		today,
	).First(&session).Error

	// ✅ 今天没开始过 → 创建
	if errors.Is(err, gorm.ErrRecordNotFound) {
		session = ledger2.LedgerSession{
			BotID:       l.botModel.BotID,
			ChatGroupID: l.chatGroupID,
			WorkDate:    today,
			IsActive:    1,
		}

		if err := global.GVA_MYSQL.Create(&session).Error; err != nil {
			return err
		}

		return l.reply("✅ 今日记账已开始")
	} else if err != nil {
		return err
	}

	// ✅ 今天已经开始过
	if session.IsActive == 1 {
		return l.reply("⚠️ 今日已经开始记账了")
	}

	// ✅ 如果被关过 → 重新开启
	session.IsActive = 1

	if err := global.GVA_MYSQL.Save(&session).Error; err != nil {
		return err
	}

	return l.reply("✅ 今日记账已重新开启")
}

func (l *StartLedgerHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}
	msg := botapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}

func IsTodayActive(botID, chatGroupID int64) (bool, error) {
	var session ledger2.LedgerSession

	today := time.Now().Format("2006-01-02")

	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ? AND work_date = ? AND is_active = 1",
		botID,
		chatGroupID,
		today,
	).First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
