
package bot

import (
	"context"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
    botReq "github.com/msean/botmanager/server/model/bot/request"
)

type BanRecordService struct {}
// CreateBanRecord 创建封禁记录记录
// Author [yourname](https://github.com/yourname)
func (banRecordService *BanRecordService) CreateBanRecord(ctx context.Context, banRecord *bot.BanRecord) (err error) {
	err = global.GVA_DB.Create(banRecord).Error
	return err
}

// DeleteBanRecord 删除封禁记录记录
// Author [yourname](https://github.com/yourname)
func (banRecordService *BanRecordService)DeleteBanRecord(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&bot.BanRecord{},"id = ?",ID).Error
	return err
}

// DeleteBanRecordByIds 批量删除封禁记录记录
// Author [yourname](https://github.com/yourname)
func (banRecordService *BanRecordService)DeleteBanRecordByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]bot.BanRecord{},"id in ?",IDs).Error
	return err
}

// UpdateBanRecord 更新封禁记录记录
// Author [yourname](https://github.com/yourname)
func (banRecordService *BanRecordService)UpdateBanRecord(ctx context.Context, banRecord bot.BanRecord) (err error) {
	err = global.GVA_DB.Model(&bot.BanRecord{}).Where("id = ?",banRecord.ID).Updates(&banRecord).Error
	return err
}

// GetBanRecord 根据ID获取封禁记录记录
// Author [yourname](https://github.com/yourname)
func (banRecordService *BanRecordService)GetBanRecord(ctx context.Context, ID string) (banRecord bot.BanRecord, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&banRecord).Error
	return
}
// GetBanRecordInfoList 分页获取封禁记录记录
// Author [yourname](https://github.com/yourname)
func (banRecordService *BanRecordService)GetBanRecordInfoList(ctx context.Context, info botReq.BanRecordSearch) (list []bot.BanRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&bot.BanRecord{})
    var banRecords []bot.BanRecord
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

	err = db.Find(&banRecords).Error
	return  banRecords, total, err
}
func (banRecordService *BanRecordService)GetBanRecordPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
