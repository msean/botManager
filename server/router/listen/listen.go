package listen

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type ListenRoter struct{}

func (s *ListenRoter) InitListenRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	botMgrRouter := Router.Group("listen").Use(middleware.OperationRecord())
	{
		botMgrRouter.GET("choice", listenApi.Choice)
		botMgrRouter.GET("query", listenApi.Query)
		botMgrRouter.POST("export", listenApi.Export)
	}
}
