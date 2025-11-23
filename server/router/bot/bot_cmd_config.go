package bot

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type BotCmdConfigRouter struct {}

// InitBotCmdConfigRouter 初始化 机器人命令配置 路由信息
func (s *BotCmdConfigRouter) InitBotCmdConfigRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	botCmdConfigRouter := Router.Group("botCmdConfig").Use(middleware.OperationRecord())
	botCmdConfigRouterWithoutRecord := Router.Group("botCmdConfig")
	botCmdConfigRouterWithoutAuth := PublicRouter.Group("botCmdConfig")
	{
		botCmdConfigRouter.POST("createBotCmdConfig", botCmdConfigApi.CreateBotCmdConfig)   // 新建机器人命令配置
		botCmdConfigRouter.DELETE("deleteBotCmdConfig", botCmdConfigApi.DeleteBotCmdConfig) // 删除机器人命令配置
		botCmdConfigRouter.DELETE("deleteBotCmdConfigByIds", botCmdConfigApi.DeleteBotCmdConfigByIds) // 批量删除机器人命令配置
		botCmdConfigRouter.PUT("updateBotCmdConfig", botCmdConfigApi.UpdateBotCmdConfig)    // 更新机器人命令配置
	}
	{
		botCmdConfigRouterWithoutRecord.GET("findBotCmdConfig", botCmdConfigApi.FindBotCmdConfig)        // 根据ID获取机器人命令配置
		botCmdConfigRouterWithoutRecord.GET("getBotCmdConfigList", botCmdConfigApi.GetBotCmdConfigList)  // 获取机器人命令配置列表
	}
	{
	    botCmdConfigRouterWithoutAuth.GET("getBotCmdConfigPublic", botCmdConfigApi.GetBotCmdConfigPublic)  // 机器人命令配置开放接口
	}
}
