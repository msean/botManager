package bot

import (
	"context"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
)

type BotService struct{}

// CreateBot 创建机器人记录
// Author [yourname](https://github.com/yourname)
func (bot_mgrService *BotService) CreateBot(ctx context.Context, bot_mgr *bot.Bot) (err error) {
	err = global.GVA_DB.Create(bot_mgr).Error
	return err
}

// DeleteBot 删除机器人记录
// Author [yourname](https://github.com/yourname)
func (bot_mgrService *BotService) DeleteBot(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&bot.Bot{}, "id = ?", ID).Error
	return err
}

// DeleteBotByIds 批量删除机器人记录
// Author [yourname](https://github.com/yourname)
func (bot_mgrService *BotService) DeleteBotByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]bot.Bot{}, "id in ?", IDs).Error
	return err
}

// UpdateBot 更新机器人记录
// Author [yourname](https://github.com/yourname)
func (bot_mgrService *BotService) UpdateBot(ctx context.Context, bot_mgr bot.Bot) (err error) {
	err = global.GVA_DB.Model(&bot.Bot{}).Where("id = ?", bot_mgr.ID).Updates(&bot_mgr).Error
	return err
}

// GetBot 根据ID获取机器人记录
// Author [yourname](https://github.com/yourname)
func (bot_mgrService *BotService) GetBot(ctx context.Context, ID string) (bot_mgr bot.Bot, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&bot_mgr).Error
	return
}

// GetBotInfoList 分页获取机器人记录
// Author [yourname](https://github.com/yourname)
func (bot_mgrService *BotService) GetBotInfoList(ctx context.Context, info botReq.BotSearch) (list []bot.Bot, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&bot.Bot{})
	var bot_mgrs []bot.Bot
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.Name != nil && *info.Name != "" {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	if info.Name != nil && *info.Token != "" {
		db = db.Where("token LIKE ?", "%"+*info.Token+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&bot_mgrs).Error
	return bot_mgrs, total, err
}
func (bot_mgrService *BotService) GetBotPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
