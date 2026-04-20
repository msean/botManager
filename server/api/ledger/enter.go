package ledger

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	LedgerApi
	LedgerPermissionApi
	LedgerAccountGroupApi
}

var (
	ledgerService             = service.ServiceGroupApp.UsageServiceGroup.LedgerService
	ledgerPermissionService   = service.ServiceGroupApp.UsageServiceGroup.LedgerPermissionService
	ledgerAccountGroupService = service.ServiceGroupApp.UsageServiceGroup.LedgerAccountGroupService
)
