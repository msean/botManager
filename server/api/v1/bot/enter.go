package bot

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	BotMsgMgrApi
	BotApi
}

var (
	bot_msg_mgrService = service.ServiceGroupApp.BotServiceGroup.BotMsgMgrService
	bot_mgrService     = service.ServiceGroupApp.BotServiceGroup.BotService
)
