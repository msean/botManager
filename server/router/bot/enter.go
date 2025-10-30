package bot

import api "github.com/msean/botmanager/server/api/v1"

type RouterGroup struct {
	BotMsgMgrRouter
	BotRouter
}

var (
	bot_msg_mgrApi = api.ApiGroupApp.BotApiGroup.BotMsgMgrApi
	bot_mgrApi     = api.ApiGroupApp.BotApiGroup.BotApi
)
