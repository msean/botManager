package tg_auto_helper

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	TgUserApi
	CollectGroupApi
}

var tgUserService = service.ServiceGroupApp.Tg_auto_helperServiceGroup.TgUserService
var collectGroupSvc = service.ServiceGroupApp.Tg_auto_helperServiceGroup.CollectGroupTaskService
