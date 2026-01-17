package chat_group

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	"github.com/msean/botmanager/server/service/cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Ledger struct {
	model       ledger.Ledger
	botModel    bot.Bot
	chatGroupID int64
}

func (l *Ledger) Match(botModel bot.Bot, update tgbotapi.Update) (match bool) {
	if update.Message == nil {
		return
	}

	msg := update.Message
	input := strings.TrimSpace(msg.Text)
	l.botModel = botModel
	l.chatGroupID = msg.Chat.ID

	if input == "" {
		return
	}

	rawInput := input

	actionType := 1
	if strings.HasPrefix(input, "+") {
		actionType = 1
		input = input[1:]
	} else if strings.HasPrefix(input, "下发") {
		actionType = 2
		input = input[1:]
	}

	sign := 1.0
	if strings.HasPrefix(input, "+") {
		input = input[1:]
	} else if strings.HasPrefix(input, "-") {
		sign = -1
		input = input[1:]
	}

	amount, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return
	}

	finalAmount := sign * amount

	_ledger := ledger.Ledger{
		OprUserID:        msg.From.ID,
		OprUserFirstName: msg.From.FirstName,
		OprUserLastName:  msg.From.LastName,
		OprUserNickname:  msg.From.UserName,

		ActionType: actionType,
		Amount:     finalAmount,

		BotID:       l.botModel.BotID,
		ChatGroupID: msg.Chat.ID,
		MessageID:   int64(msg.MessageID),
		RawInput:    rawInput,
	}

	l.model = _ledger
	return true
}

func (l *Ledger) HasPerMission() (permit bool, err error) {
	ledgerPermission := cache.NewLedgerPermissionCache(l.botModel.BotID, l.chatGroupID)
	var has bool
	if has, err = cache.CacheGetItem(ledgerPermission); err != nil {
		global.GVA_LOG.Error("Ledger HasPerMission", zap.Int64("botID", l.botModel.BotID), zap.Int64("chatGroupID", l.chatGroupID), zap.Error(err))
		return
	}
	if !has {
		global.GVA_LOG.Info("Ledger HasPerMission", zap.Int64("botID", l.botModel.BotID), zap.Int64("chatGroupID", l.chatGroupID))
		return
	}
	if !ledgerPermission.HasUserPermission(l.model.OprUserID) {
		global.GVA_LOG.Info("Ledger HasPerMission", zap.Int64("botID", l.botModel.BotID), zap.Int64("chatGroupID", l.chatGroupID), zap.Int64("userID", l.model.OprUserID))
		return
	}
	permit = true
	return
}

func (l *Ledger) Handle() (err error) {
	// 是否有权限
	var permit bool
	if permit, err = l.HasPerMission(); err != nil {
		global.GVA_LOG.Error("Ledger HasPerMission", zap.Int64("botID", l.botModel.BotID), zap.Int64("chatGroupID", l.chatGroupID), zap.Error(err))
		return
	}
	if !permit {
		global.GVA_LOG.Info("Ledger HasPerMission", zap.Int64("botID", l.botModel.BotID), zap.Int64("chatGroupID", l.chatGroupID), zap.Int64("userID", l.model.OprUserID))
		return
	}

	// 创建记账
	if err = l.Create(); err != nil {
		global.GVA_LOG.Error("Ledger Handle", zap.Error(err))
		return
	}

	var reply, url string
	if reply, url, err = l.BuildReply(global.GVA_MYSQL); err != nil {
		global.GVA_LOG.Info("Ledger HasPerMission", zap.Int64("botID", l.botModel.BotID), zap.Int64("chatGroupID", l.chatGroupID), zap.Int64("userID", l.model.OprUserID))
	}

	// var botHandler *bot_handler.Bot
	// if botHandler, err = bot_handler.NewBot(l.botModel.Token); err != nil {
	// 	global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", l.botModel.BotID), zap.Error(err))
	// 	return
	// }

	var botSender *tgbotapi.BotAPI
	if botSender, err = tgbotapi.NewBotAPI(l.botModel.Token); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(l.chatGroupID, reply)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📊 点击查看完整账单", url),
		),
	)
	if _, err = botSender.Send(msg); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", l.botModel.BotID), zap.Error(err))
	}
	// if _, err = botHandler.TgSend(l.chatGroupID, nil, reply); err != nil {
	// 	global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", l.botModel.BotID), zap.Error(err))
	// }

	return
}

func (l *Ledger) Create() error {
	return global.GVA_MYSQL.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&l.model).Error
	})
}

func (l *Ledger) BuildReply(db *gorm.DB) (content string, url string, err error) {
	bot := cache.NewBotCache(l.botModel.BotID)
	if _has, getErr := cache.CacheGetItem(bot); !_has || getErr != nil {
		global.GVA_LOG.Error("Ledger BuildReply", zap.Int64("botID", l.botModel.BotID), zap.Bool("_has", _has))
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var records []ledger.Ledger
	err = db.
		Where("chat_group_id = ?", l.chatGroupID).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Order("id asc").
		Find(&records).Error
	if err != nil {
		return
	}

	if len(records) == 0 {
		return
	}

	var (
		incomeLines []string
		payoutLines []string

		totalIncome float64
		totalPayout float64
	)

	for _, r := range records {
		line := fmt.Sprintf(
			"%s %0.2f",
			r.CreatedAt.Format("15:04:05"),
			r.Amount,
		)

		switch r.ActionType {
		case 1: // 入款
			incomeLines = append(incomeLines, line)
			totalIncome += r.Amount

		case 2: // 下发
			payoutLines = append(payoutLines, line)
			totalPayout += r.Amount
		}
	}

	shouldPayout := totalIncome
	unpaid := totalIncome - totalPayout

	url = fmt.Sprintf(
		"%s/#/ledger/full?bot_id=%d&chat_group_id=%d&idmin=%d&idmax=%d",
		global.GVA_CONFIG.System.Domain,
		l.botModel.BotID,
		l.chatGroupID,
		records[0].ID,
		records[len(records)-1].ID,
	)

	return fmt.Sprintf(
		`%s（%s）

入款（%d笔）：
%s

下发（%d笔）：
%s

总入款：%.2f
应下发：%.2f
总下发：%.2f
未下发：%.2f`,
		bot.Name,
		startOfDay.Format("2006-01-02"),
		len(incomeLines), strings.Join(incomeLines, "\n"),
		len(payoutLines), strings.Join(payoutLines, "\n"),
		totalIncome,
		shouldPayout,
		totalPayout,
		unpaid,
	), url, nil
}
