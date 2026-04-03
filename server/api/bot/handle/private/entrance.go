package private

import (
	adpublish "github.com/msean/botmanager/server/api/bot/handle/private/ad_publish"
	"github.com/msean/botmanager/server/api/bot/handle/private/ledger"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func Entrance(update botapi.Update, botModel bot.Bot) (err error) {
	if botModel.IsAdPublish == 1 {
		adpublish.Entrance(update, botModel)
	}
	if botModel.IsForLedger == 1 {
		ledger.Entrance(update, botModel)
	}
	return
}
