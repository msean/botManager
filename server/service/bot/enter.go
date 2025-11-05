package bot

type ServiceGroup struct {
	BotBanContentService
	BotService
	BanRecordService
	BotChatGroupService
	BotBanGroupMemService
	BotMsgHandlerSvc
}
