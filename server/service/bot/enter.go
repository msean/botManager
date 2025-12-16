package bot

type ServiceGroup struct {
	BotBanContentService
	BotService
	BanRecordService
	BotChatGroupService
	BotBanGroupMemService
	BotTaskService
	BotChannelService
	BotCmdConfigService
	BotHandlerSvc
	BotChatHistorySvc
}
