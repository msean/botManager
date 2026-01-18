package chat_group

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

type LedgerSetFeeRateHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	feeRate     float64
}

func (l *LedgerSetFeeRateHandler) Match(botModel bot.Bot, update tgbotapi.Update) (match bool) {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	// 必须以「费率设置」开头
	if !strings.HasPrefix(text, "费率设置") {
		return false
	}

	// 按任意空白切分
	parts := strings.Fields(text)
	if len(parts) != 2 {
		return false
	}

	fee, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return false
	}

	l.feeRate = fee
	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *LedgerSetFeeRateHandler) Handle() (err error) {
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
			// 不存在就创建
			permission = ledger.LedgerPermission{
				BotID:          l.botModel.BotID,
				ChatGroupID:    l.chatGroupID,
				CurrentFeeRate: l.feeRate,
			}
			return global.GVA_MYSQL.Create(&permission).Error
		}
		return err
	}

	// 更新费率
	return global.GVA_MYSQL.Model(&permission).
		Update("current_fee_rate", l.feeRate).Error
}
