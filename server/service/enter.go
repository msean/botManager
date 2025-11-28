package service

import (
	"time"

	"github.com/msean/botmanager/server/service/bot"
	"github.com/msean/botmanager/server/service/recharge"
	"github.com/msean/botmanager/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	BotServiceGroup      bot.ServiceGroup
	RechargeServiceGroup recharge.ServiceGroup
}

func Init() {
	bot.InitBotTaskManager()
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			recharge.CheckExpiredOrders()
		}
	}()
}
