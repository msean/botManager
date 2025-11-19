
package bot

import (
	"context"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
    botReq "github.com/msean/botmanager/server/model/bot/request"
)

type BotChannelService struct {}
// CreateBotChannel 创建机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService) CreateBotChannel(ctx context.Context, botChannel *bot.BotChannel) (err error) {
	err = global.GVA_DB.Create(botChannel).Error
	return err
}

// DeleteBotChannel 删除机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService)DeleteBotChannel(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&bot.BotChannel{},"id = ?",ID).Error
	return err
}

// DeleteBotChannelByIds 批量删除机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService)DeleteBotChannelByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]bot.BotChannel{},"id in ?",IDs).Error
	return err
}

// UpdateBotChannel 更新机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService)UpdateBotChannel(ctx context.Context, botChannel bot.BotChannel) (err error) {
	err = global.GVA_DB.Model(&bot.BotChannel{}).Where("id = ?",botChannel.ID).Updates(&botChannel).Error
	return err
}

// GetBotChannel 根据ID获取机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService)GetBotChannel(ctx context.Context, ID string) (botChannel bot.BotChannel, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&botChannel).Error
	return
}
// GetBotChannelInfoList 分页获取机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService)GetBotChannelInfoList(ctx context.Context, info botReq.BotChannelSearch) (list []bot.BotChannel, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&bot.BotChannel{})
    var botChannels []bot.BotChannel
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

	err = db.Find(&botChannels).Error
	return  botChannels, total, err
}
func (botChannelService *BotChannelService)GetBotChannelPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
