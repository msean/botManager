package bot

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type BotMsgMgrRouter struct {}

// InitBotMsgMgrRouter 初始化 机器人消息管理 路由信息
func (s *BotMsgMgrRouter) InitBotMsgMgrRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	bot_msg_mgrRouter := Router.Group("bot_msg_mgr").Use(middleware.OperationRecord())
	bot_msg_mgrRouterWithoutRecord := Router.Group("bot_msg_mgr")
	bot_msg_mgrRouterWithoutAuth := PublicRouter.Group("bot_msg_mgr")
	{
		bot_msg_mgrRouter.POST("createBotMsgMgr", bot_msg_mgrApi.CreateBotMsgMgr)   // 新建机器人消息管理
		bot_msg_mgrRouter.DELETE("deleteBotMsgMgr", bot_msg_mgrApi.DeleteBotMsgMgr) // 删除机器人消息管理
		bot_msg_mgrRouter.DELETE("deleteBotMsgMgrByIds", bot_msg_mgrApi.DeleteBotMsgMgrByIds) // 批量删除机器人消息管理
		bot_msg_mgrRouter.PUT("updateBotMsgMgr", bot_msg_mgrApi.UpdateBotMsgMgr)    // 更新机器人消息管理
	}
	{
		bot_msg_mgrRouterWithoutRecord.GET("findBotMsgMgr", bot_msg_mgrApi.FindBotMsgMgr)        // 根据ID获取机器人消息管理
		bot_msg_mgrRouterWithoutRecord.GET("getBotMsgMgrList", bot_msg_mgrApi.GetBotMsgMgrList)  // 获取机器人消息管理列表
	}
	{
	    bot_msg_mgrRouterWithoutAuth.GET("getBotMsgMgrPublic", bot_msg_mgrApi.GetBotMsgMgrPublic)  // 机器人消息管理开放接口
	}
}
