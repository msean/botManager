package dao

import (
	"strconv"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils"
	"gorm.io/gorm"
)

type rechargeDao struct{}

func newRechargeDao() *rechargeDao {
	return &rechargeDao{}
}

func (dao *rechargeDao) ExistConfig(db *gorm.DB, botID int64, publishTimes int) (has bool, err error) {
	conds := []utils.Cond{
		utils.NewWhereCond("bot_id", botID),
		utils.NewWhereCond("publish_times", publishTimes),
	}
	has, err = utils.Get(db, &recharge.RechargeConfig{}, conds...)
	return
}

// 用户最近是否有创建的订单
func (dao *rechargeDao) UserHasRecentOrder(db *gorm.DB, botID int64, userID int64) (bool, error) {

	sysCnf, loadSysErr := cache.LoadSyscnf(global.SysRepeatOrderIntervalKey, true, global.DefaultRepeatOrderInterval)
	if loadSysErr != nil {
		return false, loadSysErr
	}
	interval, _ := strconv.Atoi(sysCnf.Value)

	if interval > 0 {
		deadline := time.Now().Add(-15 * time.Minute)
		var count int64
		err := db.Model(&recharge.UserRechargeRecord{}).
			Where("bot_id = ?", botID).
			Where("user_id = ?", userID).
			Where("status = ?", global.AdRechargeCreate).
			Where("created_at >= ?", deadline).
			Count(&count).Error
		if err != nil {
			return false, err
		}

		return count > 0, nil
	}
	return false, nil
}
