package bot

import (
	"context"
	"strconv"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/service/cache"
)

type BotChannelService struct{}

// CreateBotChannel 创建机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService) CreateBotChannel(ctx context.Context, botChannel *bot.BotChannel) (err error) {
	err = global.GVA_MYSQL.Create(botChannel).Error
	return err
}

// DeleteBotChannel 删除机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService) DeleteBotChannel(ctx context.Context, ID string) (err error) {
	var id int
	if id, err = strconv.Atoi(ID); err != nil {
		return
	}
	if err = cache.ReleaseChannelModelChange(uint(id)); err != nil {
		return
	}
	return global.GVA_MYSQL.Delete(&bot.BotChannel{}, "id = ?", ID).Error
}

// DeleteBotChannelByIds 批量删除机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService) DeleteBotChannelByIds(ctx context.Context, IDs []string) (err error) {
	for _, ID := range IDs {
		var id int
		if id, _ = strconv.Atoi(ID); id > 0 {
			cache.ReleaseChannelModelChange(uint(id))
		}
	}
	err = global.GVA_MYSQL.Delete(&[]bot.BotChannel{}, "id in ?", IDs).Error
	return err
}

// UpdateBotChannel 更新机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService) UpdateBotChannel(ctx context.Context, botChannel bot.BotChannel) (err error) {
	err = global.GVA_MYSQL.Model(&bot.BotChannel{}).Where("id = ?", botChannel.ID).Updates(&botChannel).Error
	return err
}

// GetBotChannel 根据ID获取机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService) GetBotChannel(ctx context.Context, ID string) (botChannel bot.BotChannel, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&botChannel).Error
	return
}

// GetBotChannelInfoList 分页获取机器人渠道记录
// Author [yourname](https://github.com/yourname)
func (botChannelService *BotChannelService) GetBotChannelInfoList(ctx context.Context, info botReq.BotChannelSearch) (list []*bot.BotChannel, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&bot.BotChannel{})
	var botChannels []*bot.BotChannel
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Find(&botChannels).Error; err != nil {
		return
	}

	var botList []int64
	for _, object := range botChannels {
		botList = append(botList, object.BotID)
	}

	var botMapper map[int64]bot.Bot
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList); err != nil {
		return
	}
	for _, object := range botChannels {
		object.BotName = botMapper[object.BotID].Name
	}
	return botChannels, total, err
}
func (botChannelService *BotChannelService) GetBotChannelPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
