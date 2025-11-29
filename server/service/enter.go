package service

import (
	"fmt"
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
	// 每一分钟去检查过期订单
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("ReconcileAccounts panic:", r)
					}
					recharge.CheckExpiredOrders()
				}()
			}()
		}
	}()

	// 每25s 检查收款记录
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("ReconcileAccounts panic:", r)
					}
				}()
				recharge.ReconcileAccounts()
			}()
			time.Sleep(20 * time.Second)
		}
	}()
}
