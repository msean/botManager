package ledger2

import (
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func Dispatch(botModel bot.Bot, update botapi.Update, chatGroup cache.BotChatGroupCache) {
	if update.CallbackQuery != nil {
		handleCallback(botModel, update)
		return
	} else {
		Handle(botModel, update)
	}
}
