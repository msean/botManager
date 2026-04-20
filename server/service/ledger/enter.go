package ledger

type ServiceGroup struct {
	LedgerService
	LedgerPermissionService
	LedgerAccountGroupService
}
