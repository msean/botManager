package bot

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type BotRouter struct {}

// InitBotRouter 初始化 机器人 路由信息
func (s *BotRouter) InitBotRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	bot_mgrRouter := Router.Group("bot_mgr").Use(middleware.OperationRecord())
	bot_mgrRouterWithoutRecord := Router.Group("bot_mgr")
	bot_mgrRouterWithoutAuth := PublicRouter.Group("bot_mgr")
	{
		bot_mgrRouter.POST("createBot", bot_mgrApi.CreateBot)   // 新建机器人
		bot_mgrRouter.DELETE("deleteBot", bot_mgrApi.DeleteBot) // 删除机器人
		bot_mgrRouter.DELETE("deleteBotByIds", bot_mgrApi.DeleteBotByIds) // 批量删除机器人
		bot_mgrRouter.PUT("updateBot", bot_mgrApi.UpdateBot)    // 更新机器人
	}
	{
		bot_mgrRouterWithoutRecord.GET("findBot", bot_mgrApi.FindBot)        // 根据ID获取机器人
		bot_mgrRouterWithoutRecord.GET("getBotList", bot_mgrApi.GetBotList)  // 获取机器人列表
	}
	{
	    bot_mgrRouterWithoutAuth.GET("getBotPublic", bot_mgrApi.GetBotPublic)  // 机器人开放接口
	}
}
