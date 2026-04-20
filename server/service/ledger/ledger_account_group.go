
package ledger

import (
	"context"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/ledger"
    ledgerReq "github.com/msean/botmanager/server/model/ledger/request"
)

type LedgerAccountGroupService struct {}
// CreateLedgerAccountGroup 创建记账账号组记录
// Author [yourname](https://github.com/yourname)
func (ledgerAccountGroupService *LedgerAccountGroupService) CreateLedgerAccountGroup(ctx context.Context, ledgerAccountGroup *ledger.LedgerAccountGroup) (err error) {
	err = global.GVA_MYSQL.Create(ledgerAccountGroup).Error
	return err
}

// DeleteLedgerAccountGroup 删除记账账号组记录
// Author [yourname](https://github.com/yourname)
func (ledgerAccountGroupService *LedgerAccountGroupService)DeleteLedgerAccountGroup(ctx context.Context, ID string) (err error) {
	err = global.GVA_MYSQL.Delete(&ledger.LedgerAccountGroup{},"id = ?",ID).Error
	return err
}

// DeleteLedgerAccountGroupByIds 批量删除记账账号组记录
// Author [yourname](https://github.com/yourname)
func (ledgerAccountGroupService *LedgerAccountGroupService)DeleteLedgerAccountGroupByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]ledger.LedgerAccountGroup{},"id in ?",IDs).Error
	return err
}

// UpdateLedgerAccountGroup 更新记账账号组记录
// Author [yourname](https://github.com/yourname)
func (ledgerAccountGroupService *LedgerAccountGroupService)UpdateLedgerAccountGroup(ctx context.Context, ledgerAccountGroup ledger.LedgerAccountGroup) (err error) {
	err = global.GVA_MYSQL.Model(&ledger.LedgerAccountGroup{}).Where("id = ?",ledgerAccountGroup.ID).Updates(&ledgerAccountGroup).Error
	return err
}

// GetLedgerAccountGroup 根据ID获取记账账号组记录
// Author [yourname](https://github.com/yourname)
func (ledgerAccountGroupService *LedgerAccountGroupService)GetLedgerAccountGroup(ctx context.Context, ID string) (ledgerAccountGroup ledger.LedgerAccountGroup, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&ledgerAccountGroup).Error
	return
}
// GetLedgerAccountGroupInfoList 分页获取记账账号组记录
// Author [yourname](https://github.com/yourname)
func (ledgerAccountGroupService *LedgerAccountGroupService)GetLedgerAccountGroupInfoList(ctx context.Context, info ledgerReq.LedgerAccountGroupSearch) (list []ledger.LedgerAccountGroup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_MYSQL.Model(&ledger.LedgerAccountGroup{})
    var ledgerAccountGroups []ledger.LedgerAccountGroup
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&ledgerAccountGroups).Error
	return  ledgerAccountGroups, total, err
}
func (ledgerAccountGroupService *LedgerAccountGroupService)GetLedgerAccountGroupPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
