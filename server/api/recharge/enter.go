package recharge

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	RechargeConfigApi
	UserRechargeRecordApi
	UserWalletApi
}

var (
	rechargeConfigService     = service.ServiceGroupApp.RechargeServiceGroup.RechargeConfigService
	userRechargeRecordService = service.ServiceGroupApp.RechargeServiceGroup.UserRechargeRecordService
	userWalletService         = service.ServiceGroupApp.RechargeServiceGroup.UserWalletService
)
