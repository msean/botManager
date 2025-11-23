package recharge

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type RechargeConfigRouter struct {}

// InitRechargeConfigRouter 初始化 充值配置 路由信息
func (s *RechargeConfigRouter) InitRechargeConfigRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	rechargeConfigRouter := Router.Group("rechargeConfig").Use(middleware.OperationRecord())
	rechargeConfigRouterWithoutRecord := Router.Group("rechargeConfig")
	rechargeConfigRouterWithoutAuth := PublicRouter.Group("rechargeConfig")
	{
		rechargeConfigRouter.POST("createRechargeConfig", rechargeConfigApi.CreateRechargeConfig)   // 新建充值配置
		rechargeConfigRouter.DELETE("deleteRechargeConfig", rechargeConfigApi.DeleteRechargeConfig) // 删除充值配置
		rechargeConfigRouter.DELETE("deleteRechargeConfigByIds", rechargeConfigApi.DeleteRechargeConfigByIds) // 批量删除充值配置
		rechargeConfigRouter.PUT("updateRechargeConfig", rechargeConfigApi.UpdateRechargeConfig)    // 更新充值配置
	}
	{
		rechargeConfigRouterWithoutRecord.GET("findRechargeConfig", rechargeConfigApi.FindRechargeConfig)        // 根据ID获取充值配置
		rechargeConfigRouterWithoutRecord.GET("getRechargeConfigList", rechargeConfigApi.GetRechargeConfigList)  // 获取充值配置列表
	}
	{
	    rechargeConfigRouterWithoutAuth.GET("getRechargeConfigPublic", rechargeConfigApi.GetRechargeConfigPublic)  // 充值配置开放接口
	}
}
