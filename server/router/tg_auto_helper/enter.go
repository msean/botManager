package tg_auto_helper

import api "github.com/msean/botmanager/server/api"

type RouterGroup struct {
	TgUserRouter
	CollectGroupRouter
}

var tgUserApi = api.ApiGroupApp.TgAutoHelperGroup.TgUserApi
var collectGroupApi = api.ApiGroupApp.TgAutoHelperGroup.CollectGroupApi
