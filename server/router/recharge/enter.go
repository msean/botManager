package recharge

import (
	api "github.com/msean/botmanager/server/api"
)

type RouterGroup struct {
	RechargeConfigRouter
	UserRechargeRecordRouter
	UserWalletRouter
	AdPublishRecordRouter
}

var (
	rechargeConfigApi     = api.ApiGroupApp.RechargeApiGroup.RechargeConfigApi
	userRechargeRecordApi = api.ApiGroupApp.RechargeApiGroup.UserRechargeRecordApi
	userWalletApi         = api.ApiGroupApp.RechargeApiGroup.UserWalletApi
	adPublishRecordApi    = api.ApiGroupApp.RechargeApiGroup.AdPublishRecordApi
)
