package bot

import (
	api "github.com/msean/botmanager/server/api"
)

type RouterGroup struct {
	BotBanContentRouter
	BotRouter
	BanRecordRouter
	BotChatGroupRouter
	BotBanGroupMemRouter
	BotTaskRouter
	BotChannelRouter
	BotCmdConfigRouter
	BotMsgMassRouter
}

var (
	botBanContentApi  = api.ApiGroupApp.BotApiGroup.BotBanContentApi
	botApi            = api.ApiGroupApp.BotApiGroup.BotApi
	botMsgHandler     = api.ApiGroupApp.BotApiGroup.BotMsgHandler
	banRecordApi      = api.ApiGroupApp.BotApiGroup.BanRecordApi
	botChatGroupApi   = api.ApiGroupApp.BotApiGroup.BotChatGroupApi
	botBanGroupMemApi = api.ApiGroupApp.BotApiGroup.BotBanGroupMemApi
	taskApi           = api.ApiGroupApp.BotApiGroup.BotTaskApi
	botChannelApi     = api.ApiGroupApp.BotApiGroup.BotChannelApi
	botCmdConfigApi   = api.ApiGroupApp.BotApiGroup.BotCmdConfigApi
	botMsgMassApi     = api.ApiGroupApp.BotApiGroup.BotMsgMassApi
)
