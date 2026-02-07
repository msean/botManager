package tg_auto_helper

import api "github.com/msean/botmanager/server/api"

type RouterGroup struct{ TgUserRouter }

var tgUserApi = api.ApiGroupApp.TgAutoHelperGroup.TgUserApi
