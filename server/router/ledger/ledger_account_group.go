package ledger

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type LedgerAccountGroupRouter struct{}

// InitLedgerAccountGroupRouter 初始化 记账账号组 路由信息
func (s *LedgerAccountGroupRouter) InitLedgerAccountGroupRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	ledgerAccountGroupRouter := Router.Group("ledgerAccountGroup").Use(middleware.OperationRecord())
	ledgerAccountGroupRouterWithoutRecord := Router.Group("ledgerAccountGroup")
	// ledgerAccountGroupRouterWithoutAuth := PublicRouter.Group("ledgerAccountGroup")
	{
		ledgerAccountGroupRouter.POST("createLedgerAccountGroup", ledgerAccountGroupApi.CreateLedgerAccountGroup)   // 新建记账账号组
		ledgerAccountGroupRouter.DELETE("deleteLedgerAccountGroup", ledgerAccountGroupApi.DeleteLedgerAccountGroup) // 删除记账账号组
		ledgerAccountGroupRouter.PUT("updateLedgerAccountGroup", ledgerAccountGroupApi.UpdateLedgerAccountGroup)    // 更新记账账号组
	}
	{
		ledgerAccountGroupRouterWithoutRecord.GET("findLedgerAccountGroup", ledgerAccountGroupApi.FindLedgerAccountGroup)       // 根据ID获取记账账号组
		ledgerAccountGroupRouterWithoutRecord.GET("getLedgerAccountGroupList", ledgerAccountGroupApi.GetLedgerAccountGroupList) // 获取记账账号组列表
	}
}
