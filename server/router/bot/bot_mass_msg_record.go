package bot

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type BotMassMsgRecordRouter struct {}

// InitBotMassMsgRecordRouter 初始化 群发历史记录 路由信息
func (s *BotMassMsgRecordRouter) InitBotMassMsgRecordRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	botMassMsgRecordRouter := Router.Group("botMassMsgRecord").Use(middleware.OperationRecord())
	botMassMsgRecordRouterWithoutRecord := Router.Group("botMassMsgRecord")
	botMassMsgRecordRouterWithoutAuth := PublicRouter.Group("botMassMsgRecord")
	{
		botMassMsgRecordRouter.POST("createBotMassMsgRecord", botMassMsgRecordApi.CreateBotMassMsgRecord)   // 新建群发历史记录
		botMassMsgRecordRouter.DELETE("deleteBotMassMsgRecord", botMassMsgRecordApi.DeleteBotMassMsgRecord) // 删除群发历史记录
		botMassMsgRecordRouter.DELETE("deleteBotMassMsgRecordByIds", botMassMsgRecordApi.DeleteBotMassMsgRecordByIds) // 批量删除群发历史记录
		botMassMsgRecordRouter.PUT("updateBotMassMsgRecord", botMassMsgRecordApi.UpdateBotMassMsgRecord)    // 更新群发历史记录
	}
	{
		botMassMsgRecordRouterWithoutRecord.GET("findBotMassMsgRecord", botMassMsgRecordApi.FindBotMassMsgRecord)        // 根据ID获取群发历史记录
		botMassMsgRecordRouterWithoutRecord.GET("getBotMassMsgRecordList", botMassMsgRecordApi.GetBotMassMsgRecordList)  // 获取群发历史记录列表
	}
	{
	    botMassMsgRecordRouterWithoutAuth.GET("getBotMassMsgRecordPublic", botMassMsgRecordApi.GetBotMassMsgRecordPublic)  // 群发历史记录开放接口
	}
}
