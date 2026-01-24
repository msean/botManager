package ledger

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type LedgerPermissionRouter struct {}

// InitLedgerPermissionRouter 初始化 帐薄权限管理 路由信息
func (s *LedgerPermissionRouter) InitLedgerPermissionRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	ledgerPermissionRouter := Router.Group("ledgerPermission").Use(middleware.OperationRecord())
	ledgerPermissionRouterWithoutRecord := Router.Group("ledgerPermission")
	ledgerPermissionRouterWithoutAuth := PublicRouter.Group("ledgerPermission")
	{
		ledgerPermissionRouter.POST("createLedgerPermission", ledgerPermissionApi.CreateLedgerPermission)   // 新建帐薄权限管理
		ledgerPermissionRouter.DELETE("deleteLedgerPermission", ledgerPermissionApi.DeleteLedgerPermission) // 删除帐薄权限管理
		ledgerPermissionRouter.DELETE("deleteLedgerPermissionByIds", ledgerPermissionApi.DeleteLedgerPermissionByIds) // 批量删除帐薄权限管理
		ledgerPermissionRouter.PUT("updateLedgerPermission", ledgerPermissionApi.UpdateLedgerPermission)    // 更新帐薄权限管理
	}
	{
		ledgerPermissionRouterWithoutRecord.GET("findLedgerPermission", ledgerPermissionApi.FindLedgerPermission)        // 根据ID获取帐薄权限管理
		ledgerPermissionRouterWithoutRecord.GET("getLedgerPermissionList", ledgerPermissionApi.GetLedgerPermissionList)  // 获取帐薄权限管理列表
	}
	{
	    ledgerPermissionRouterWithoutAuth.GET("getLedgerPermissionPublic", ledgerPermissionApi.GetLedgerPermissionPublic)  // 帐薄权限管理开放接口
	}
}
