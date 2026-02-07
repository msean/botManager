
package bot

import (
	"context"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
    botReq "github.com/msean/botmanager/server/model/bot/request"
)

type BotMassMsgRecordService struct {}
// CreateBotMassMsgRecord 创建群发历史记录记录
// Author [yourname](https://github.com/yourname)
func (botMassMsgRecordService *BotMassMsgRecordService) CreateBotMassMsgRecord(ctx context.Context, botMassMsgRecord *bot.BotMassMsgRecord) (err error) {
	err = global.GVA_MYSQL.Create(botMassMsgRecord).Error
	return err
}

// DeleteBotMassMsgRecord 删除群发历史记录记录
// Author [yourname](https://github.com/yourname)
func (botMassMsgRecordService *BotMassMsgRecordService)DeleteBotMassMsgRecord(ctx context.Context, ID string) (err error) {
	err = global.GVA_MYSQL.Delete(&bot.BotMassMsgRecord{},"id = ?",ID).Error
	return err
}

// DeleteBotMassMsgRecordByIds 批量删除群发历史记录记录
// Author [yourname](https://github.com/yourname)
func (botMassMsgRecordService *BotMassMsgRecordService)DeleteBotMassMsgRecordByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]bot.BotMassMsgRecord{},"id in ?",IDs).Error
	return err
}

// UpdateBotMassMsgRecord 更新群发历史记录记录
// Author [yourname](https://github.com/yourname)
func (botMassMsgRecordService *BotMassMsgRecordService)UpdateBotMassMsgRecord(ctx context.Context, botMassMsgRecord bot.BotMassMsgRecord) (err error) {
	err = global.GVA_MYSQL.Model(&bot.BotMassMsgRecord{}).Where("id = ?",botMassMsgRecord.ID).Updates(&botMassMsgRecord).Error
	return err
}

// GetBotMassMsgRecord 根据ID获取群发历史记录记录
// Author [yourname](https://github.com/yourname)
func (botMassMsgRecordService *BotMassMsgRecordService)GetBotMassMsgRecord(ctx context.Context, ID string) (botMassMsgRecord bot.BotMassMsgRecord, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&botMassMsgRecord).Error
	return
}
// GetBotMassMsgRecordInfoList 分页获取群发历史记录记录
// Author [yourname](https://github.com/yourname)
func (botMassMsgRecordService *BotMassMsgRecordService)GetBotMassMsgRecordInfoList(ctx context.Context, info botReq.BotMassMsgRecordSearch) (list []bot.BotMassMsgRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_MYSQL.Model(&bot.BotMassMsgRecord{})
    var botMassMsgRecords []bot.BotMassMsgRecord
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

	err = db.Find(&botMassMsgRecords).Error
	return  botMassMsgRecords, total, err
}
func (botMassMsgRecordService *BotMassMsgRecordService)GetBotMassMsgRecordPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
