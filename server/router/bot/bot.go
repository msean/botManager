package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type BotRouter struct{}

// InitBotRouter 初始化 机器人 路由信息
func (s *BotRouter) InitBotRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	botMgrRouter := Router.Group("bot_mgr").Use(middleware.OperationRecord())
	botMgrRouterWithoutRecord := Router.Group("bot_mgr")
	botMgrRouterWithoutAuth := PublicRouter.Group("bot_mgr")
	{
		botMgrRouter.POST("create", botApi.CreateBot)               // 新建机器人
		botMgrRouter.DELETE("delete", botApi.DeleteBot)             // 删除机器人
		botMgrRouter.DELETE("delete_by_ids", botApi.DeleteBotByIds) // 批量删除机器人
		botMgrRouter.PUT("update", botApi.UpdateBot)                // 更新机器人
		botMgrRouter.POST("unban_user", botApi.UnbanUser)           // 更新机器人
	}
	{
		botMgrRouterWithoutRecord.GET("get", botApi.FindBot)     // 根据ID获取机器人
		botMgrRouterWithoutRecord.GET("list", botApi.GetBotList) // 获取机器人列表
		botMgrRouterWithoutRecord.GET("choice", botApi.All)      //
		botMgrRouterWithoutRecord.GET("choice_with_chat_group", botApi.AllWithChatGroupAndChannel)
	}
	{
		botMgrRouterWithoutAuth.GET("public_get", botApi.GetBotPublic) // 机器人开放接口
	}

	botHandlerRouter := PublicRouter.Group("bot")
	botHandlerRouter.POST("/webhook/:botUUID", botMsgHandler.Handle)
}
