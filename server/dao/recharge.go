package dao

import (
	"errors"
	"strconv"
	"time"

	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	sysCnf, loadSysErr := cache.LoadSyscnf(constant.SysRepeatOrderIntervalKey, true, constant.DefaultRepeatOrderInterval)
	if loadSysErr != nil {
		return false, loadSysErr
	}
	interval, _ := strconv.Atoi(sysCnf.Value)

	if interval > 0 {
		var count int64
		err := db.Model(&recharge.UserRechargeRecord{}).
			Where("bot_id = ?", botID).
			Where("user_id = ?", userID).
			Where("status = ?", constant.AdRechargeCreate).
			Where("created_at >= ?", time.Now().Add(-constant.OrderLeftPaid*time.Minute)).
			Count(&count).Error
		if err != nil {
			return false, err
		}

		return count > 0, nil
	}
	return false, nil
}

// func (dao *rechargeDao) CancelOrder(db *gorm.DB, botID int64, userID int64, updateID int64) error {
// 	return db.Model(&recharge.UserRechargeRecord{}).
// 		Where("bot_id = ? AND user_id = ? AND update_id = ?", botID, userID, updateID).
// 		Update("status", constant.AdRechargeCancel).Error
// }

func (dao *rechargeDao) GetUserWallet(db *gorm.DB, botID, userID int64, userName string) (wallet recharge.UserWallet, err error) {
	err = db.
		Where("bot_id = ? AND user_id = ?", botID, userID).
		Attrs(recharge.UserWallet{
			UserName: userName,
		}).
		FirstOrCreate(&wallet, recharge.UserWallet{
			UserID: userID,
			BotID:  botID,
		}).Error
	return
}

func (dao *rechargeDao) AddBalance(db *gorm.DB, botID, userID int64, amount float64) (balance float64, err error) {

	err = db.Transaction(func(tx *gorm.DB) error {

		var wallet recharge.UserWallet

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("bot_id = ? AND user_id = ?", botID, userID).
			First(&wallet).Error; err != nil {
			return err
		}

		wallet.Balance += amount

		if err := tx.Model(&wallet).
			Update("balance", wallet.Balance).Error; err != nil {
			return err
		}

		balance = wallet.Balance
		return nil
	})

	return
}

func (dao *rechargeDao) ReduceBalance(db *gorm.DB, botID, userID int64, amount float64) (balance float64, err error) {

	err = db.Transaction(func(tx *gorm.DB) error {

		var wallet recharge.UserWallet

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("bot_id = ? AND user_id = ?", botID, userID).
			First(&wallet).Error; err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("余额不足")
		}

		wallet.Balance -= amount

		if err := tx.Model(&wallet).
			Update("balance", wallet.Balance).Error; err != nil {
			return err
		}

		balance = wallet.Balance
		return nil
	})

	return
}
