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
		botRouter.InitBotMsgMgrRouter(privateGroup, publicGroup)
		botRouter.InitBotRouter(privateGroup, publicGroup)
	}
}
