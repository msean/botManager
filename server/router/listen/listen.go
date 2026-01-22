package listen

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type ListenRoter struct{}

// InitBotRouter 初始化 机器人 路由信息
func (s *ListenRoter) InitBotRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	botMgrRouter := Router.Group("listen").Use(middleware.OperationRecord())
	// botMgrRouterWithoutRecord := Router.Group("listen")
	// botMgrRouterWithoutAuth := PublicRouter.Group("listen")
	{
		botMgrRouter.GET("create", botApi.CreateBot) // 新建机器人
	}

	botHandlerRouter := PublicRouter.Group("bot")
	botHandlerRouter.POST("/webhook/:botUUID", botMsgHandler.Handle)
}
