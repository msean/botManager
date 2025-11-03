package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type BotBanContentRouter struct{}

// InitBotBanContentRouter 初始化 机器人消息管理 路由信息
func (s *BotBanContentRouter) InitBotBanContentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	botBanContentRouter := Router.Group("bot_ban_content").Use(middleware.OperationRecord())
	BotBanContentRouterWithoutRecord := Router.Group("bot_ban_content")
	BotBanContentRouterWithoutAuth := PublicRouter.Group("bot_ban_content")
	{
		botBanContentRouter.POST("create", botBanContentApi.CreateBotBanContent)               // 新建机器人消息管理
		botBanContentRouter.DELETE("delete", botBanContentApi.DeleteBotBanContent)             // 删除机器人消息管理
		botBanContentRouter.DELETE("delete_by_ids", botBanContentApi.DeleteBotBanContentByIds) // 批量删除机器人消息管理
		botBanContentRouter.PUT("update", botBanContentApi.UpdateBotBanContent)                // 更新机器人消息管理
	}
	{
		BotBanContentRouterWithoutRecord.GET("get", botBanContentApi.FindBotBanContent)     // 根据ID获取机器人消息管理
		BotBanContentRouterWithoutRecord.GET("list", botBanContentApi.GetBotBanContentList) // 获取机器人消息管理列表
	}
	{
		BotBanContentRouterWithoutAuth.GET("public_get", botBanContentApi.GetBotBanContentPublic) // 机器人消息管理开放接口
	}
}
