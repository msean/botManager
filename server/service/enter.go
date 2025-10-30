package service

import (
	"github.com/msean/botmanager/server/service/bot"
	"github.com/msean/botmanager/server/service/example"
	"github.com/msean/botmanager/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	BotServiceGroup     bot.ServiceGroup
}
