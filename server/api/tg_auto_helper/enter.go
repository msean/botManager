package tg_auto_helper

import "github.com/msean/botmanager/server/service"

type ApiGroup struct{ TgUserApi }

var tgUserService = service.ServiceGroupApp.Tg_auto_helperServiceGroup.TgUserService
