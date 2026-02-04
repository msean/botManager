package rate_ledger

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type CurrencyRateSetHandler struct {
	botModel     bot.Bot
	chatGroupID  int64
	currencyRate float64
	ShouldPermissionAwareWithOutAdmin
}

func (l *CurrencyRateSetHandler) Match(botModel bot.Bot, update botapi.Update) (match bool) {
	text := strings.TrimSpace(update.Message.Text)

	if !strings.HasPrefix(text, "设置汇率") {
		return false
	}

	parts := strings.Fields(text)
	if len(parts) != 2 {
		return false
	}

	fee, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return false
	}

	l.currencyRate = fee
	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *CurrencyRateSetHandler) Handle() (err error) {
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
				BotID:                  l.botModel.BotID,
				ChatGroupID:            l.chatGroupID,
				CurrentCurrencyFeeRate: l.currencyRate,
			}
			return global.GVA_MYSQL.Create(&permission).Error
		}
		return err
	}

	if err = global.GVA_MYSQL.Model(&permission).
		Update("current_fee_rate", l.currencyRate).Error; err != nil {
		return
	}
	if e := cache.NewLedgerPermissionCache(l.botModel.BotID, l.chatGroupID).Release(); e != nil {
		global.GVA_LOG.Error("Ledger parse amount failed",
			zap.Int64("botID", l.botModel.BotID),
			zap.Int64("chatGroupID", l.chatGroupID),
			zap.Error(e),
		)
	}
	l.reply(fmt.Sprintf("汇率设置成功 当前费率为%s", formatFeeRate(l.currencyRate)))
	return
}

func (l *CurrencyRateSetHandler) reply(text string) (err error) {
	var botSender *botapi.BotAPI
	botSender, err = botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return
	}
	msg := botapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
