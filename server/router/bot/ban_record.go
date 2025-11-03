package bot

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type BanRecordRouter struct {}

// InitBanRecordRouter 初始化 封禁记录 路由信息
func (s *BanRecordRouter) InitBanRecordRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	banRecordRouter := Router.Group("banRecord").Use(middleware.OperationRecord())
	banRecordRouterWithoutRecord := Router.Group("banRecord")
	banRecordRouterWithoutAuth := PublicRouter.Group("banRecord")
	{
		banRecordRouter.POST("createBanRecord", banRecordApi.CreateBanRecord)   // 新建封禁记录
		banRecordRouter.DELETE("deleteBanRecord", banRecordApi.DeleteBanRecord) // 删除封禁记录
		banRecordRouter.DELETE("deleteBanRecordByIds", banRecordApi.DeleteBanRecordByIds) // 批量删除封禁记录
		banRecordRouter.PUT("updateBanRecord", banRecordApi.UpdateBanRecord)    // 更新封禁记录
	}
	{
		banRecordRouterWithoutRecord.GET("findBanRecord", banRecordApi.FindBanRecord)        // 根据ID获取封禁记录
		banRecordRouterWithoutRecord.GET("getBanRecordList", banRecordApi.GetBanRecordList)  // 获取封禁记录列表
	}
	{
	    banRecordRouterWithoutAuth.GET("getBanRecordPublic", banRecordApi.GetBanRecordPublic)  // 封禁记录开放接口
	}
}
