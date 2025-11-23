package bot

import "github.com/msean/botmanager/server/service"

var (
	botChatGroupService   = service.ServiceGroupApp.BotServiceGroup.BotChatGroupService
	botBanGroupMemService = service.ServiceGroupApp.BotServiceGroup.BotBanGroupMemService
	taskService           = service.ServiceGroupApp.BotServiceGroup.BotTaskService
	botChannelService     = service.ServiceGroupApp.BotServiceGroup.BotChannelService
	botCmdConfigService   = service.ServiceGroupApp.BotServiceGroup.BotCmdConfigService
)
