package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type BotChatGroupRouter struct{}

// InitBotChatGroupRouter 初始化 机器人群组列表 路由信息
func (s *BotChatGroupRouter) InitBotChatGroupRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	botChatGroupRouter := Router.Group("botChatGroup").Use(middleware.OperationRecord())
	botChatGroupRouterWithoutRecord := Router.Group("botChatGroup")
	botChatGroupRouterWithoutAuth := PublicRouter.Group("botChatGroup")
	{
		botChatGroupRouter.POST("createBotChatGroup", botChatGroupApi.CreateBotChatGroup)             // 新建机器人群组列表
		botChatGroupRouter.DELETE("deleteBotChatGroup", botChatGroupApi.DeleteBotChatGroup)           // 删除机器人群组列表
		botChatGroupRouter.DELETE("deleteBotChatGroupByIds", botChatGroupApi.DeleteBotChatGroupByIds) // 批量删除机器人群组列表
		botChatGroupRouter.PUT("updateBotChatGroup", botChatGroupApi.UpdateBotChatGroup)              // 更新机器人群组列表
		botChatGroupRouter.GET("getBotChatGroupClassifyList", botChatGroupApi.GetClassfyList)
		botChatGroupRouter.POST("saveBotChatGroupClassify", botChatGroupApi.SaveClassify)
		botChatGroupRouter.DELETE("deleteBotChatGroupClassify", botChatGroupApi.DeleteClassify)
		botChatGroupRouter.GET("chooseChatGroupClassify", botChatGroupApi.ClassifyChoice)
	}
	{
		botChatGroupRouterWithoutRecord.GET("findBotChatGroup", botChatGroupApi.FindBotChatGroup)       // 根据ID获取机器人群组列表
		botChatGroupRouterWithoutRecord.GET("getBotChatGroupList", botChatGroupApi.GetBotChatGroupList) // 获取机器人群组列表列表
		botChatGroupRouterWithoutRecord.GET("chatHistory", botChatGroupApi.ChatHistory)
	}
	{
		botChatGroupRouterWithoutAuth.GET("getBotChatGroupPublic", botChatGroupApi.GetBotChatGroupPublic) // 机器人群组列表开放接口
	}
}
