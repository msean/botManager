package recharge

import (
	"context"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/recharge"
	rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
)

type AdPublishRecordService struct{}

// CreateAdPublishRecord 创建广告发布记录记录
// Author [yourname](https://github.com/yourname)
func (adPublishRecordService *AdPublishRecordService) CreateAdPublishRecord(ctx context.Context, adPublishRecord *recharge.AdPublishRecord) (err error) {
	err = global.GVA_MYSQL.Create(adPublishRecord).Error
	return err
}

// DeleteAdPublishRecord 删除广告发布记录记录
// Author [yourname](https://github.com/yourname)
func (adPublishRecordService *AdPublishRecordService) DeleteAdPublishRecord(ctx context.Context, ID string) (err error) {
	err = global.GVA_MYSQL.Delete(&recharge.AdPublishRecord{}, "id = ?", ID).Error
	return err
}

// DeleteAdPublishRecordByIds 批量删除广告发布记录记录
// Author [yourname](https://github.com/yourname)
func (adPublishRecordService *AdPublishRecordService) DeleteAdPublishRecordByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]recharge.AdPublishRecord{}, "id in ?", IDs).Error
	return err
}

// UpdateAdPublishRecord 更新广告发布记录记录
// Author [yourname](https://github.com/yourname)
func (adPublishRecordService *AdPublishRecordService) UpdateAdPublishRecord(ctx context.Context, adPublishRecord recharge.AdPublishRecord) (err error) {
	err = global.GVA_MYSQL.Model(&recharge.AdPublishRecord{}).Where("id = ?", adPublishRecord.ID).Updates(&adPublishRecord).Error
	return err
}

// GetAdPublishRecord 根据ID获取广告发布记录记录
// Author [yourname](https://github.com/yourname)
func (adPublishRecordService *AdPublishRecordService) GetAdPublishRecord(ctx context.Context, ID string) (adPublishRecord recharge.AdPublishRecord, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&adPublishRecord).Error
	return
}

// GetAdPublishRecordInfoList 分页获取广告发布记录记录
// Author [yourname](https://github.com/yourname)
func (adPublishRecordService *AdPublishRecordService) GetAdPublishRecordInfoList(ctx context.Context, info rechargeReq.AdPublishRecordSearch) (list []*recharge.AdPublishRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&recharge.AdPublishRecord{})
	var adPublishRecords []*recharge.AdPublishRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.UserID != 0 {
		db = db.Where("user_id = ?", info.UserID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Find(&adPublishRecords).Error; err != nil {
		return
	}

	var botList, channelList []int64
	for _, object := range adPublishRecords {
		botList = append(botList, object.BotID)
		channelList = append(channelList, object.ChannelID)
	}

	var botMapper map[int64]bot.Bot
	var channelMapper map[int64]bot.BotChannel
	botMapper, _ = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList)
	channelMapper, _ = dao.BotChannelDao.MappByChannelIDList(global.GVA_MYSQL, channelList)

	for _, object := range adPublishRecords {
		object.BotName = botMapper[object.BotID].Name
		object.ChannelName = channelMapper[object.ChannelID].ChannelName
	}
	return adPublishRecords, total, err
}

func (adPublishRecordService *AdPublishRecordService) GetAdPublishRecordPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
