package rate_ledger

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type CurrencyRateShowHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	ShouldPermissionAwareWithOutAdmin
}

func (l *CurrencyRateShowHandler) Match(botModel bot.Bot, update botapi.Update) (match bool) {

	text := strings.TrimSpace(update.Message.Text)

	if text != "查看汇率" {
		return false
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *CurrencyRateShowHandler) Handle() (err error) {
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
			return l.reply("当前群聊尚未设置汇率")
		}
		return err
	}

	return l.reply(
		"当前汇率为：" + formatFeeRate(permission.CurrentCurrencyFeeRate),
	)
}

func (l *CurrencyRateShowHandler) reply(text string) (err error) {
	var botSender *botapi.BotAPI
	botSender, err = botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return
	}
	msg := botapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
