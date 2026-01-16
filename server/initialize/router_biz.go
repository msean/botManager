package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/router"
)

func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}
func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]
	holder(publicGroup, privateGroup)
	{
		botRouter := router.RouterGroupApp.Bot
		botRouter.InitBotBanContentRouter(privateGroup, publicGroup)
		botRouter.InitBotRouter(privateGroup, publicGroup)
		botRouter.InitBanRecordRouter(privateGroup, publicGroup)
		botRouter.InitBotChatGroupRouter(privateGroup, publicGroup)
		botRouter.InitBotBanGroupMemRouter(privateGroup, publicGroup)
		botRouter.InitBotTaskRouter(privateGroup, publicGroup)
		publicRouter := router.RouterGroupApp.Public
		publicRouter.InitMedioRouter(privateGroup, publicGroup)
		botRouter.InitBotChannelRouter(privateGroup, publicGroup)
		botRouter.InitBotCmdConfigRouter(privateGroup, publicGroup)
	}
	{
		rechargeRouter := router.RouterGroupApp.Recharge
		rechargeRouter.InitRechargeConfigRouter(privateGroup, publicGroup)
		rechargeRouter.InitUserRechargeRecordRouter(privateGroup, publicGroup)
		rechargeRouter.InitUserWalletRouter(privateGroup, publicGroup)
		rechargeRouter.InitAdPublishRecordRouter(privateGroup, publicGroup)
	}
	{
		usageRouter := router.RouterGroupApp.Usage
		usageRouter.InitLedgerRouter(privateGroup, publicGroup)
		usageRouter.InitLedgerPermissionRouter(privateGroup, publicGroup)
	}
}
