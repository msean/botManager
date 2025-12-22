package recharge

import (
	"context"
	"errors"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/recharge"
	rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
)

type UserWalletService struct{}

// CreateUserWallet 创建用户钱包记录
// Author [yourname](https://github.com/yourname)
func (userWalletService *UserWalletService) CreateUserWallet(ctx context.Context, userWallet *recharge.UserWallet) (err error) {
	var has bool
	if has, err = dao.RechargeDao.WalletExist(global.GVA_MYSQL, userWallet.BotID, userWallet.UserID); err != nil {
		return
	}
	if has {
		err = errors.New("该机器人下存在用户的钱包记录")
		return
	}
	err = global.GVA_MYSQL.Create(userWallet).Error
	return err
}

// DeleteUserWallet 删除用户钱包记录
// Author [yourname](https://github.com/yourname)
func (userWalletService *UserWalletService) DeleteUserWallet(ctx context.Context, ID string) (err error) {
	err = global.GVA_MYSQL.Delete(&recharge.UserWallet{}, "id = ?", ID).Error
	return err
}

// DeleteUserWalletByIds 批量删除用户钱包记录
// Author [yourname](https://github.com/yourname)
func (userWalletService *UserWalletService) DeleteUserWalletByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]recharge.UserWallet{}, "id in ?", IDs).Error
	return err
}

// UpdateUserWallet 更新用户钱包记录
// Author [yourname](https://github.com/yourname)
func (userWalletService *UserWalletService) UpdateUserWallet(ctx context.Context, userWallet recharge.UserWallet) (err error) {
	err = global.GVA_MYSQL.Model(&recharge.UserWallet{}).Where("id = ?", userWallet.ID).Updates(&userWallet).Error
	return err
}

// GetUserWallet 根据ID获取用户钱包记录
// Author [yourname](https://github.com/yourname)
func (userWalletService *UserWalletService) GetUserWallet(ctx context.Context, ID string) (userWallet recharge.UserWallet, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&userWallet).Error
	return
}

// GetUserWalletInfoList 分页获取用户钱包记录
// Author [yourname](https://github.com/yourname)
func (userWalletService *UserWalletService) GetUserWalletInfoList(ctx context.Context, info rechargeReq.UserWalletSearch) (list []*recharge.UserWallet, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&recharge.UserWallet{})
	var userWallets []*recharge.UserWallet
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.UserID != 0 {
		db = db.Where("user_id=?", info.UserID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Find(&userWallets).Error; err != nil {
		return
	}

	var botList []int64
	for _, object := range userWallets {
		botList = append(botList, object.BotID)
	}

	var botMapper map[int64]bot.Bot
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList); err != nil {
		return
	}
	for _, object := range userWallets {
		object.BotName = botMapper[object.BotID].Name
	}
	return userWallets, total, err
}
func (userWalletService *UserWalletService) GetUserWalletPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
