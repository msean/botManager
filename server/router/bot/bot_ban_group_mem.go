package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type BotBanGroupMemRouter struct{}

// InitBotBanGroupMemRouter 初始化 封禁成员设置 路由信息
func (s *BotBanGroupMemRouter) InitBotBanGroupMemRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	botBanGroupMemRouter := Router.Group("botBanGroupMem").Use(middleware.OperationRecord())
	botBanGroupMemRouterWithoutRecord := Router.Group("botBanGroupMem")
	botBanGroupMemRouterWithoutAuth := PublicRouter.Group("botBanGroupMem")
	{
		botBanGroupMemRouter.POST("createBotBanGroupMem", botBanGroupMemApi.CreateBotBanGroupMem)             // 新建封禁成员设置
		botBanGroupMemRouter.DELETE("deleteBotBanGroupMem", botBanGroupMemApi.DeleteBotBanGroupMem)           // 删除封禁成员设置
		botBanGroupMemRouter.DELETE("deleteBotBanGroupMemByIds", botBanGroupMemApi.DeleteBotBanGroupMemByIds) // 批量删除封禁成员设置
		botBanGroupMemRouter.PUT("updateBotBanGroupMem", botBanGroupMemApi.UpdateBotBanGroupMem)              // 更新封禁成员设置
	}
	{
		botBanGroupMemRouterWithoutRecord.GET("findBotBanGroupMem", botBanGroupMemApi.FindBotBanGroupMem)       // 根据ID获取封禁成员设置
		botBanGroupMemRouterWithoutRecord.GET("getBotBanGroupMemList", botBanGroupMemApi.GetBotBanGroupMemList) // 获取封禁成员设置列表
	}
	{
		botBanGroupMemRouterWithoutAuth.GET("getBotBanGroupMemPublic", botBanGroupMemApi.GetBotBanGroupMemPublic) // 封禁成员设置开放接口
	}
}
