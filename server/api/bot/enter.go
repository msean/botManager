package bot

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	BotBanContentApi
	BotApi
	BotMsgHandler
	BanRecordApi
}

var (
	BotBanContentService = service.ServiceGroupApp.BotServiceGroup.BotBanContentService
	botMgrService        = service.ServiceGroupApp.BotServiceGroup.BotService
	banRecordService     = service.ServiceGroupApp.BotServiceGroup.BanRecordService
)
