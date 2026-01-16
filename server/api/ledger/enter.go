package ledger

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	LedgerApi
	LedgerPermissionApi
}

var (
	ledgerService           = service.ServiceGroupApp.UsageServiceGroup.LedgerService
	ledgerPermissionService = service.ServiceGroupApp.UsageServiceGroup.LedgerPermissionService
)
