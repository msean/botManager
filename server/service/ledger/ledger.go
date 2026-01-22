package usage

import (
	"context"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	ledgerReq "github.com/msean/botmanager/server/model/ledger/request"
)

type LedgerService struct{}

// CreateLedger 创建帐薄记录
// Author [yourname](https://github.com/yourname)
func (ledgerService *LedgerService) CreateLedger(ctx context.Context, ledger *ledger.Ledger) (err error) {
	err = global.GVA_MYSQL.Create(ledger).Error
	return err
}

// DeleteLedger 删除帐薄记录
// Author [yourname](https://github.com/yourname)
func (ledgerService *LedgerService) DeleteLedger(ctx context.Context, ID string) (err error) {
	err = global.GVA_MYSQL.Delete(&ledger.Ledger{}, "id = ?", ID).Error
	return err
}

// DeleteLedgerByIds 批量删除帐薄记录
// Author [yourname](https://github.com/yourname)
func (ledgerService *LedgerService) DeleteLedgerByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]ledger.Ledger{}, "id in ?", IDs).Error
	return err
}

// UpdateLedger 更新帐薄记录
// Author [yourname](https://github.com/yourname)
func (ledgerService *LedgerService) UpdateLedger(ctx context.Context, _ledger ledger.Ledger) (err error) {
	err = global.GVA_MYSQL.Model(&ledger.Ledger{}).Where("id = ?", _ledger.ID).Updates(&_ledger).Error
	return err
}

// GetLedger 根据ID获取帐薄记录
// Author [yourname](https://github.com/yourname)
func (ledgerService *LedgerService) GetLedger(ctx context.Context, ID string) (_ledger ledger.Ledger, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&_ledger).Error
	return
}

// GetLedgerInfoList 分页获取帐薄记录
// Author [yourname](https://github.com/yourname)
func (ledgerService *LedgerService) GetLedgerInfoList(ctx context.Context, info ledgerReq.LedgerSearch) (list []*ledger.Ledger, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&ledger.Ledger{})
	var ledgers []*ledger.Ledger
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

	if err = db.Find(&ledgers).Error; err != nil {
		return
	}

	var botList []int64
	var chatGroupList []int64
	for _, object := range ledgers {
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

	for _, object := range ledgers {
		object.BotName = botMapper[object.BotID].Name
		object.ChatGroupName = chatGroupMapper[object.ChatGroupID].ChatGroupName
	}
	return ledgers, total, err
}

func (ledgerService *LedgerService) GetLedgerPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
