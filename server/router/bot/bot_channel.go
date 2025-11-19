package bot

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type BotChannelRouter struct {}

// InitBotChannelRouter 初始化 机器人渠道 路由信息
func (s *BotChannelRouter) InitBotChannelRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	botChannelRouter := Router.Group("botChannel").Use(middleware.OperationRecord())
	botChannelRouterWithoutRecord := Router.Group("botChannel")
	botChannelRouterWithoutAuth := PublicRouter.Group("botChannel")
	{
		botChannelRouter.POST("createBotChannel", botChannelApi.CreateBotChannel)   // 新建机器人渠道
		botChannelRouter.DELETE("deleteBotChannel", botChannelApi.DeleteBotChannel) // 删除机器人渠道
		botChannelRouter.DELETE("deleteBotChannelByIds", botChannelApi.DeleteBotChannelByIds) // 批量删除机器人渠道
		botChannelRouter.PUT("updateBotChannel", botChannelApi.UpdateBotChannel)    // 更新机器人渠道
	}
	{
		botChannelRouterWithoutRecord.GET("findBotChannel", botChannelApi.FindBotChannel)        // 根据ID获取机器人渠道
		botChannelRouterWithoutRecord.GET("getBotChannelList", botChannelApi.GetBotChannelList)  // 获取机器人渠道列表
	}
	{
	    botChannelRouterWithoutAuth.GET("getBotChannelPublic", botChannelApi.GetBotChannelPublic)  // 机器人渠道开放接口
	}
}
