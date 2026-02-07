package bot

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	BotBanContentApi
	BotApi
	BotMsgHandler
	BanRecordApi
	BotBanGroupMemApi
	BotChatGroupApi
	BotTaskApi
	BotChannelApi
	BotCmdConfigApi
	BotMsgMassApi
	BotMassMsgRecordApi
}

var (
	BotBanContentService    = service.ServiceGroupApp.BotServiceGroup.BotBanContentService
	botMgrService           = service.ServiceGroupApp.BotServiceGroup.BotService
	banRecordService        = service.ServiceGroupApp.BotServiceGroup.BanRecordService
	botChatGroupService     = service.ServiceGroupApp.BotServiceGroup.BotChatGroupService
	botBanGroupMemService   = service.ServiceGroupApp.BotServiceGroup.BotBanGroupMemService
	taskService             = service.ServiceGroupApp.BotServiceGroup.BotTaskService
	botChannelService       = service.ServiceGroupApp.BotServiceGroup.BotChannelService
	botCmdConfigService     = service.ServiceGroupApp.BotServiceGroup.BotCmdConfigService
	botMsgMassService       = service.ServiceGroupApp.BotServiceGroup.BotMsgMassService
	botMassMsgRecordService = service.ServiceGroupApp.BotServiceGroup.BotMassMsgRecordService
)
