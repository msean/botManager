package service

import (
	"fmt"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/service/bot"
	"github.com/msean/botmanager/server/service/recharge"
	"github.com/msean/botmanager/server/service/system"
	"time"
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
		for range // 每一分钟去检查过期订单
		// 每25s 检查收款记录
		ticker.C {
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
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("ReconcileAccounts panic:", r)
					}
				}()
				if !global.GVA_CONFIG.System.UnMatchPayment {
					recharge.ReconcileAccounts()
				}
			}()
			time.Sleep(15 * time.Second)
		}
	}()
}
