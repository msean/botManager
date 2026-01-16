package usage

import (
	"context"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	ledgerReq "github.com/msean/botmanager/server/model/ledger/request"
)

type LedgerPermissionService struct{}

// CreateLedgerPermission 创建帐薄权限管理记录
// Author [yourname](https://github.com/yourname)
func (ledgerPermissionService *LedgerPermissionService) CreateLedgerPermission(ctx context.Context, ledgerPermission *ledger.LedgerPermission) (err error) {
	err = global.GVA_MYSQL.Create(ledgerPermission).Error
	return err
}

// DeleteLedgerPermission 删除帐薄权限管理记录
// Author [yourname](https://github.com/yourname)
func (ledgerPermissionService *LedgerPermissionService) DeleteLedgerPermission(ctx context.Context, ID string) (err error) {
	err = global.GVA_MYSQL.Delete(&ledger.LedgerPermission{}, "id = ?", ID).Error
	return err
}

// DeleteLedgerPermissionByIds 批量删除帐薄权限管理记录
// Author [yourname](https://github.com/yourname)
func (ledgerPermissionService *LedgerPermissionService) DeleteLedgerPermissionByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]ledger.LedgerPermission{}, "id in ?", IDs).Error
	return err
}

// UpdateLedgerPermission 更新帐薄权限管理记录
// Author [yourname](https://github.com/yourname)
func (ledgerPermissionService *LedgerPermissionService) UpdateLedgerPermission(ctx context.Context, ledgerPermission ledger.LedgerPermission) (err error) {
	err = global.GVA_MYSQL.Model(&ledger.LedgerPermission{}).Where("id = ?", ledgerPermission.ID).Updates(&ledgerPermission).Error
	return err
}

// GetLedgerPermission 根据ID获取帐薄权限管理记录
// Author [yourname](https://github.com/yourname)
func (ledgerPermissionService *LedgerPermissionService) GetLedgerPermission(ctx context.Context, ID string) (ledgerPermission ledger.LedgerPermission, err error) {
	if err = global.GVA_MYSQL.Where("id = ?", ID).First(&ledgerPermission).Error; err != nil {
		return
	}

	bot, _, _ := dao.BotDao.FromBotID(global.GVA_MYSQL, int(ledgerPermission.BotID))
	botChatGroup, _, _ := dao.BotChatGroupDao.FromID(global.GVA_MYSQL, int(ledgerPermission.ChatGroupID))

	ledgerPermission.BotName = bot.Name
	ledgerPermission.ChatGroupName = botChatGroup.ChatGroupName
	return
}

// GetLedgerPermissionInfoList 分页获取帐薄权限管理记录
// Author [yourname](https://github.com/yourname)
func (ledgerPermissionService *LedgerPermissionService) GetLedgerPermissionInfoList(ctx context.Context, info ledgerReq.LedgerPermissionSearch) (list []*ledger.LedgerPermission, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&ledger.LedgerPermission{})
	var ledgerPermissions []*ledger.LedgerPermission
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Find(&ledgerPermissions).Error; err != nil {
		return
	}

	var botList []int64
	var chatGroupList []int64
	for _, object := range ledgerPermissions {
		botList = append(botList, object.BotID)
		chatGroupList = append(chatGroupList, object.ChatGroupID)
	}

	var botMapper map[int64]bot.Bot
	var chatGroupMapper map[int64]bot.BotChatGroup
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList); err != nil {
		return
	}

	if chatGroupMapper, err = dao.BotChatGroupDao.MappByChatGroupIDList(global.GVA_MYSQL, chatGroupList); err != nil {
		return
	}

	for _, object := range ledgerPermissions {
		object.BotName = botMapper[object.BotID].Name
		object.ChatGroupName = chatGroupMapper[object.ChatGroupID].ChatGroupName
	}
	return ledgerPermissions, total, err
}
func (ledgerPermissionService *LedgerPermissionService) GetLedgerPermissionPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
