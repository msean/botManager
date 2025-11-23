package bot

import "github.com/msean/botmanager/server/service/bot/handle"

type ServiceGroup struct {
	BotBanContentService
	BotService
	BanRecordService
	BotChatGroupService
	BotBanGroupMemService
	handle.BotMsgHandlerSvc
	BotTaskService
	BotChannelService
	BotCmdConfigService
}
