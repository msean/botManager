package private

import (
	adpublish "github.com/msean/botmanager/server/api/bot/handle/private/ad_publish"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func Entrance(update botapi.Update, botModel bot.Bot) (err error) {
	if botModel.IsAdPublish == 1 {
		adpublish.Handle(update, botModel.Token, int64(botModel.BotID))
	}
	return
}
