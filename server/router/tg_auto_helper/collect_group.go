package tg_auto_helper

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type CollectGroupRouter struct{}

func (s *CollectGroupRouter) Init(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	router := Router.Group("collect_group").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("collect_group")
	{
		router.POST("create", collectGroupApi.Create)
		router.DELETE("delete", collectGroupApi.Delete)
	}
	{
		routerWithoutRecord.GET("list", collectGroupApi.List)
		routerWithoutRecord.GET("list_collect_group_info", collectGroupApi.ListCollectGroupInfo)
	}
}
