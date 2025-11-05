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
}

var (
	botBanContentApi  = api.ApiGroupApp.BotApiGroup.BotBanContentApi
	botMsgApi         = api.ApiGroupApp.BotApiGroup.BotApi
	botMsgHandler     = api.ApiGroupApp.BotApiGroup.BotMsgHandler
	banRecordApi      = api.ApiGroupApp.BotApiGroup.BanRecordApi
	botChatGroupApi   = api.ApiGroupApp.BotApiGroup.BotChatGroupApi
	botBanGroupMemApi = api.ApiGroupApp.BotApiGroup.BotBanGroupMemApi
)
