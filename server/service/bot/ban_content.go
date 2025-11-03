package bot

import (
	"context"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
)

type BotBanContentService struct{}

// CreateBotBanContent 创建机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) CreateBotBanContent(ctx context.Context, bot_msg_mgr *bot.BotBanContent) (err error) {
	err = global.GVA_DB.Create(bot_msg_mgr).Error
	return err
}

// DeleteBotBanContent 删除机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) DeleteBotBanContent(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&bot.BotBanContent{}, "id = ?", ID).Error
	return err
}

// DeleteBotBanContentByIds 批量删除机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) DeleteBotBanContentByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]bot.BotBanContent{}, "id in ?", IDs).Error
	return err
}

// UpdateBotBanContent 更新机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) UpdateBotBanContent(ctx context.Context, bot_msg_mgr bot.BotBanContent) (err error) {
	err = global.GVA_DB.Model(&bot.BotBanContent{}).Where("id = ?", bot_msg_mgr.ID).Updates(&bot_msg_mgr).Error
	return err
}

// GetBotBanContent 根据ID获取机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) GetBotBanContent(ctx context.Context, ID string) (bot_msg_mgr bot.BotBanContent, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&bot_msg_mgr).Error
	return
}

// GetBotBanContentInfoList 分页获取机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) GetBotBanContentInfoList(ctx context.Context, info botReq.BotBanContentSearch) (list []*bot.BotBanContent, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&bot.BotBanContent{})
	var botBanContents []*bot.BotBanContent
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.BotID != 0 {
		db = db.Where("bot_id =?", info.BotID)
	}

	if info.BanContent != "" {
		db = db.Where("ban_content LIKE ?", "%"+info.BanContent+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Find(&botBanContents).Error; err != nil {
		return
	}
	var botList []int
	for _, botBanContent := range botBanContents {
		botList = append(botList, int(botBanContent.BotID))
	}

	var botMapper map[int]bot.Bot
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_DB, botList); err != nil {
		return
	}

	for _, botBanContent := range botBanContents {
		botBanContent.BotName = botMapper[int(botBanContent.BotID)].Name
	}
	return botBanContents, total, err
}
func (svc *BotBanContentService) GetBotBanContentPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
