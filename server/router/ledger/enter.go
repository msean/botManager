package ledger

import (
	api "github.com/msean/botmanager/server/api"
)

type RouterGroup struct {
	LedgerRouter
	LedgerPermissionRouter
	LedgerAccountGroupRouter
}

var (
	ledgerApi             = api.ApiGroupApp.LedgerApiGroup.LedgerApi
	ledgerPermissionApi   = api.ApiGroupApp.LedgerApiGroup.LedgerPermissionApi
	ledgerAccountGroupApi = api.ApiGroupApp.LedgerApiGroup.LedgerAccountGroupApi
)
