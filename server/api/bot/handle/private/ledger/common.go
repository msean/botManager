package ledger

import (
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func replyText(token string, chatGroupID int64, text string) (err error) {
	var botSender *botapi.BotAPI
	botSender, err = botapi.NewBotAPI(token)
	if err != nil {
		return
	}
	msg := botapi.NewMessage(chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
