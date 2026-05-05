package ledger

import (
	"github.com/msean/botmanager/server/global"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func replyText(token string, chatGroupID int64, text string) (err error) {
	var botSender *botapi.BotAPI
	botSender, err = botapi.NewBotAPI(token)
	if err != nil {
		return
	}
	msg := botapi.NewMessage(chatGroupID, text)
	if _, err = botSender.Send(msg); err != nil {
		global.GVA_LOG.Error("replyText", zap.Error(err))
	}
	return err
}
