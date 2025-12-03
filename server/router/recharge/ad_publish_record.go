package recharge

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type AdPublishRecordRouter struct {}

// InitAdPublishRecordRouter 初始化 广告发布记录 路由信息
func (s *AdPublishRecordRouter) InitAdPublishRecordRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	adPublishRecordRouter := Router.Group("adPublishRecord").Use(middleware.OperationRecord())
	adPublishRecordRouterWithoutRecord := Router.Group("adPublishRecord")
	adPublishRecordRouterWithoutAuth := PublicRouter.Group("adPublishRecord")
	{
		adPublishRecordRouter.POST("createAdPublishRecord", adPublishRecordApi.CreateAdPublishRecord)   // 新建广告发布记录
		adPublishRecordRouter.DELETE("deleteAdPublishRecord", adPublishRecordApi.DeleteAdPublishRecord) // 删除广告发布记录
		adPublishRecordRouter.DELETE("deleteAdPublishRecordByIds", adPublishRecordApi.DeleteAdPublishRecordByIds) // 批量删除广告发布记录
		adPublishRecordRouter.PUT("updateAdPublishRecord", adPublishRecordApi.UpdateAdPublishRecord)    // 更新广告发布记录
	}
	{
		adPublishRecordRouterWithoutRecord.GET("findAdPublishRecord", adPublishRecordApi.FindAdPublishRecord)        // 根据ID获取广告发布记录
		adPublishRecordRouterWithoutRecord.GET("getAdPublishRecordList", adPublishRecordApi.GetAdPublishRecordList)  // 获取广告发布记录列表
	}
	{
	    adPublishRecordRouterWithoutAuth.GET("getAdPublishRecordPublic", adPublishRecordApi.GetAdPublishRecordPublic)  // 广告发布记录开放接口
	}
}
