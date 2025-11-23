
package recharge

import (
	"context"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/recharge"
    rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
)

type UserRechargeRecordService struct {}
// CreateUserRechargeRecord 创建用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService) CreateUserRechargeRecord(ctx context.Context, userRechargeRecord *recharge.UserRechargeRecord) (err error) {
	err = global.GVA_DB.Create(userRechargeRecord).Error
	return err
}

// DeleteUserRechargeRecord 删除用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService)DeleteUserRechargeRecord(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&recharge.UserRechargeRecord{},"id = ?",ID).Error
	return err
}

// DeleteUserRechargeRecordByIds 批量删除用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService)DeleteUserRechargeRecordByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]recharge.UserRechargeRecord{},"id in ?",IDs).Error
	return err
}

// UpdateUserRechargeRecord 更新用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService)UpdateUserRechargeRecord(ctx context.Context, userRechargeRecord recharge.UserRechargeRecord) (err error) {
	err = global.GVA_DB.Model(&recharge.UserRechargeRecord{}).Where("id = ?",userRechargeRecord.ID).Updates(&userRechargeRecord).Error
	return err
}

// GetUserRechargeRecord 根据ID获取用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService)GetUserRechargeRecord(ctx context.Context, ID string) (userRechargeRecord recharge.UserRechargeRecord, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&userRechargeRecord).Error
	return
}
// GetUserRechargeRecordInfoList 分页获取用户充值记录记录
// Author [yourname](https://github.com/yourname)
func (userRechargeRecordService *UserRechargeRecordService)GetUserRechargeRecordInfoList(ctx context.Context, info rechargeReq.UserRechargeRecordSearch) (list []recharge.UserRechargeRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&recharge.UserRechargeRecord{})
    var userRechargeRecords []recharge.UserRechargeRecord
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

	err = db.Find(&userRechargeRecords).Error
	return  userRechargeRecords, total, err
}
func (userRechargeRecordService *UserRechargeRecordService)GetUserRechargeRecordPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
