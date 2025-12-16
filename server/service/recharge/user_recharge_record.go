package recharge

import (
	"context"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/recharge"
	rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
)

type UserRechargeRecordService struct{}

// CreateUserRechargeRecord 创建用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService) CreateUserRechargeRecord(ctx context.Context, userRechargeRecord *recharge.UserRechargeRecord) (err error) {
	err = global.GVA_MYSQL.Create(userRechargeRecord).Error
	return err
}

// DeleteUserRechargeRecord 删除用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService) DeleteUserRechargeRecord(ctx context.Context, ID string) (err error) {
	err = global.GVA_MYSQL.Delete(&recharge.UserRechargeRecord{}, "id = ?", ID).Error
	return err
}

// DeleteUserRechargeRecordByIds 批量删除用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService) DeleteUserRechargeRecordByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]recharge.UserRechargeRecord{}, "id in ?", IDs).Error
	return err
}

// UpdateUserRechargeRecord 更新用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService) UpdateUserRechargeRecord(ctx context.Context, userRechargeRecord recharge.UserRechargeRecord) (err error) {
	err = global.GVA_MYSQL.Model(&recharge.UserRechargeRecord{}).Where("id = ?", userRechargeRecord.ID).Updates(&userRechargeRecord).Error
	return err
}

// GetUserRechargeRecord 根据ID获取用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService) GetUserRechargeRecord(ctx context.Context, ID string) (userRechargeRecord recharge.UserRechargeRecord, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&userRechargeRecord).Error
	return
}

// GetUserRechargeRecordInfoList 分页获取用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService) GetUserRechargeRecordInfoList(ctx context.Context, info rechargeReq.UserRechargeRecordSearch) (list []*recharge.UserRechargeRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&recharge.UserRechargeRecord{})
	var userRechargeRecords []*recharge.UserRechargeRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if info.BotID != 0 {
		db = db.Where("bot_id=?", info.BotID)
	}
	if info.Status != 0 {
		db = db.Where("status=?", info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	db = db.Order("id desc")
	if err = db.Find(&userRechargeRecords).Error; err != nil {
		return
	}

	var botList []int64
	for _, object := range userRechargeRecords {
		botList = append(botList, object.BotID)
	}

	var botMapper map[int64]bot.Bot
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList); err != nil {
		return
	}

	for _, object := range userRechargeRecords {
		object.BotName = botMapper[object.BotID].Name
	}
	return userRechargeRecords, total, err
}
func (userRechargeRecordService *UserRechargeRecordService) GetUserRechargeRecordPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
