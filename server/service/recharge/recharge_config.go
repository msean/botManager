package recharge

import (
	"context"
	"fmt"
	"strconv"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/recharge"
	rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
	"github.com/msean/botmanager/server/service/cache"
)

type RechargeConfigService struct{}

// CreateRechargeConfig 创建充值配置记录
// Author [yourname](https://github.com/yourname)
func (rechargeConfigService *RechargeConfigService) CreateRechargeConfig(ctx context.Context, rechargeConfig *recharge.RechargeConfig) (err error) {
	var has bool
	has, err = dao.RechargeDao.ExistConfig(global.GVA_DB, rechargeConfig.BotID, int(rechargeConfig.PublishTimes))
	if err != nil {
		return
	}
	if has {
		return fmt.Errorf("存在该机器人发布次数%d的配置", rechargeConfig.PublishTimes)
	}
	err = global.GVA_DB.Create(rechargeConfig).Error
	return err
}

// DeleteRechargeConfig 删除充值配置记录
// Author [yourname](https://github.com/yourname)
func (rechargeConfigService *RechargeConfigService) DeleteRechargeConfig(ctx context.Context, ID string) (err error) {
	var id int
	if id, err = strconv.Atoi(ID); err != nil {
		return
	}
	cache.ReleaseRechargeCnf(id)
	if err = global.GVA_DB.Delete(&recharge.RechargeConfig{}, "id = ?", ID).Error; err != nil {
		return
	}
	return err
}

// DeleteRechargeConfigByIds 批量删除充值配置记录
// Author [yourname](https://github.com/yourname)
func (rechargeConfigService *RechargeConfigService) DeleteRechargeConfigByIds(ctx context.Context, IDs []string) (err error) {
	for _, _id := range IDs {
		var id int
		if id, err = strconv.Atoi(_id); err == nil {
			cache.ReleaseRechargeCnf(id)
		}
	}
	return global.GVA_DB.Delete(&[]recharge.RechargeConfig{}, "id in ?", IDs).Error
}

// UpdateRechargeConfig 更新充值配置记录
// Author [yourname](https://github.com/yourname)
func (rechargeConfigService *RechargeConfigService) UpdateRechargeConfig(ctx context.Context, rechargeConfig recharge.RechargeConfig) (err error) {
	cache.ReleaseRechargeCnf(int(rechargeConfig.ID))
	err = global.GVA_DB.Model(&recharge.RechargeConfig{}).Where("id = ?", rechargeConfig.ID).Updates(&rechargeConfig).Error
	return err
}

// GetRechargeConfig 根据ID获取充值配置记录
// Author [yourname](https://github.com/yourname)
func (rechargeConfigService *RechargeConfigService) GetRechargeConfig(ctx context.Context, ID string) (rechargeConfig recharge.RechargeConfig, err error) {
	if err = global.GVA_DB.Where("id = ?", ID).First(&rechargeConfig).Error; err != nil {
		return
	}
	return
}

// GetRechargeConfigInfoList 分页获取充值配置记录
// Author [yourname](https://github.com/yourname)
func (rechargeConfigService *RechargeConfigService) GetRechargeConfigInfoList(ctx context.Context, info rechargeReq.RechargeConfigSearch) (list []*recharge.RechargeConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&recharge.RechargeConfig{})
	var rechargeConfigs []*recharge.RechargeConfig
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.BotID != 0 {
		db = db.Where("bot_id=?", info.BotID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Find(&rechargeConfigs).Error; err != nil {
		return
	}

	var botList []int
	for _, object := range rechargeConfigs {
		botList = append(botList, int(object.BotID))
	}

	var botMapper map[int]bot.Bot
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_DB, botList); err != nil {
		return
	}

	for _, object := range rechargeConfigs {
		object.BotName = botMapper[int(object.BotID)].Name
	}
	return rechargeConfigs, total, err
}

func (rechargeConfigService *RechargeConfigService) GetRechargeConfigPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
