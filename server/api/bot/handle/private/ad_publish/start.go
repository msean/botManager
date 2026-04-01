package adpublish

import (
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func StartHandlerfunc(update botapi.Update, token string, cfg cache.BotCmdCache) {
	SendCfgMessage(update.Message.Chat.ID, token, cfg, constant.ButtonTypeKeyBoard)
}
