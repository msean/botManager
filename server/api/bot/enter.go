package bot

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	BotBanContentApi
	BotApi
	BotMsgHandler
	BanRecordApi
	BotBanGroupMemApi
	BotChatGroupApi
}

var (
	BotBanContentService  = service.ServiceGroupApp.BotServiceGroup.BotBanContentService
	botMgrService         = service.ServiceGroupApp.BotServiceGroup.BotService
	banRecordService      = service.ServiceGroupApp.BotServiceGroup.BanRecordService
	botChatGroupService   = service.ServiceGroupApp.BotServiceGroup.BotChatGroupService
	botBanGroupMemService = service.ServiceGroupApp.BotServiceGroup.BotBanGroupMemService
	botMsgHandlerSvc      = service.ServiceGroupApp.BotServiceGroup.BotMsgHandlerSvc
)
