package bot

import "github.com/msean/botmanager/server/service"

var (
	botChatGroupService   = service.ServiceGroupApp.BotServiceGroup.BotChatGroupService
	botBanGroupMemService = service.ServiceGroupApp.BotServiceGroup.BotBanGroupMemService
)
