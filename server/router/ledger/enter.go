package ledger

import (
	api "github.com/msean/botmanager/server/api"
)

type RouterGroup struct {
	LedgerRouter
	LedgerPermissionRouter
}

var (
	ledgerApi           = api.ApiGroupApp.LedgerApiGroup.LedgerApi
	ledgerPermissionApi = api.ApiGroupApp.LedgerApiGroup.LedgerPermissionApi
)
