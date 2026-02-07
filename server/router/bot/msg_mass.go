package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type BotMsgMassRouter struct{}

// InitBotMsgMassRouter 初始化 机器人群发 路由信息
func (s *BotMsgMassRouter) InitBotMsgMassRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	botMsgMassRouter := Router.Group("botMsgMass").Use(middleware.OperationRecord())
	botMsgMassRouterWithoutRecord := Router.Group("botMsgMass")
	botMsgMassRouterWithoutAuth := PublicRouter.Group("botMsgMass")
	{
		botMsgMassRouter.POST("createBotMsgMass", botMsgMassApi.CreateBotMsgMass)             // 新建机器人群发
		botMsgMassRouter.DELETE("deleteBotMsgMass", botMsgMassApi.DeleteBotMsgMass)           // 删除机器人群发
		botMsgMassRouter.DELETE("deleteBotMsgMassByIds", botMsgMassApi.DeleteBotMsgMassByIds) // 批量删除机器人群发
		botMsgMassRouter.PUT("updateBotMsgMass", botMsgMassApi.UpdateBotMsgMass)              // 更新机器人群发
		botMsgMassRouter.POST("sendBotMsgMass", botMsgMassApi.SendBotMsgMass)                 // 更新机器人群发
	}
	{
		botMsgMassRouterWithoutRecord.GET("findBotMsgMass", botMsgMassApi.FindBotMsgMass)       // 根据ID获取机器人群发
		botMsgMassRouterWithoutRecord.GET("getBotMsgMassList", botMsgMassApi.GetBotMsgMassList) // 获取机器人群发列表
	}
	{
		botMsgMassRouterWithoutAuth.GET("getBotMsgMassPublic", botMsgMassApi.GetBotMsgMassPublic) // 机器人群发开放接口
	}
}
