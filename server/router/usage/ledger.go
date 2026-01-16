package usage

import (
	"github.com/gin-gonic/gin"
)

type LedgerRouter struct{}

// InitLedgerRouter 初始化 帐薄 路由信息
func (s *LedgerRouter) InitLedgerRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	// ledgerRouter := Router.Group("ledger").Use(middleware.OperationRecord())
	ledgerRouterWithoutRecord := Router.Group("ledger")
	// ledgerRouterWithoutAuth := PublicRouter.Group("ledger")
	// {
	// 	ledgerRouter.POST("createLedger", ledgerApi.CreateLedger)   // 新建帐薄
	// 	ledgerRouter.DELETE("deleteLedger", ledgerApi.DeleteLedger) // 删除帐薄
	// 	ledgerRouter.DELETE("deleteLedgerByIds", ledgerApi.DeleteLedgerByIds) // 批量删除帐薄
	// 	ledgerRouter.PUT("updateLedger", ledgerApi.UpdateLedger)    // 更新帐薄
	// }
	{
		// ledgerRouterWithoutRecord.GET("findLedger", ledgerApi.FindLedger)       // 根据ID获取帐薄
		ledgerRouterWithoutRecord.GET("getLedgerList", ledgerApi.GetLedgerList) // 获取帐薄列表
	}
	// {
	// 	ledgerRouterWithoutAuth.GET("getLedgerPublic", ledgerApi.GetLedgerPublic) // 帐薄开放接口
	// }
}
