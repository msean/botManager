package bot

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type BotTaskRouter struct {}

// InitBotTaskRouter 初始化 任务列表 路由信息
func (s *BotTaskRouter) InitBotTaskRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	taskRouter := Router.Group("task").Use(middleware.OperationRecord())
	taskRouterWithoutRecord := Router.Group("task")
	taskRouterWithoutAuth := PublicRouter.Group("task")
	{
		taskRouter.POST("createBotTask", taskApi.CreateBotTask)   // 新建任务列表
		taskRouter.DELETE("deleteBotTask", taskApi.DeleteBotTask) // 删除任务列表
		taskRouter.DELETE("deleteBotTaskByIds", taskApi.DeleteBotTaskByIds) // 批量删除任务列表
		taskRouter.PUT("updateBotTask", taskApi.UpdateBotTask)    // 更新任务列表
	}
	{
		taskRouterWithoutRecord.GET("findBotTask", taskApi.FindBotTask)        // 根据ID获取任务列表
		taskRouterWithoutRecord.GET("getBotTaskList", taskApi.GetBotTaskList)  // 获取任务列表列表
	}
	{
	    taskRouterWithoutAuth.GET("getBotTaskPublic", taskApi.GetBotTaskPublic)  // 任务列表开放接口
	}
}
