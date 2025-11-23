package dao

import (
	"github.com/msean/botmanager/server/model/recharge"
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
