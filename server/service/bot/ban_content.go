package bot

import (
	"context"
	"strconv"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils"
	"go.uber.org/zap"
)

type BotBanContentService struct{}

// CreateBotBanContent 创建机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) CreateBotBanContent(ctx context.Context, botBanContent *bot.BotBanContent) (err error) {
	if err = global.GVA_DB.Create(botBanContent).Error; err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("id", botBanContent.ID))
		return
	}
	return err
}

// DeleteBotBanContent 删除机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) DeleteBotBanContent(ctx context.Context, ID string) (err error) {
	var id int
	if id, err = strconv.Atoi(ID); err != nil {
		return
	}
	var botContent bot.BotBanContent
	var has bool
	if has, err = utils.Get(global.GVA_DB, &botContent, utils.IDCond(ID)); !has || err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("id", id), zap.Error(err))
		return
	}
	if err = global.GVA_DB.Delete(&bot.BotBanContent{}, "id = ?", ID).Error; err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("id", id))
		return
	}

	if deleteErr := cache.NewBotBanContentListCache(botContent.BotID).Release(); deleteErr != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("BotID", botContent.BotID))
	}
	return err
}

// DeleteBotBanContentByIds 批量删除机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) DeleteBotBanContentByIds(ctx context.Context, IDs []string) (err error) {
	ids := utils.StringsToIntsIgnoreError(IDs)
	var objects []bot.BotBanContent
	if err = utils.Find(global.GVA_DB, &objects, utils.NewInCond("id", utils.IntSliceToAnySlice(ids))); err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("ids", IDs), zap.Error(err))
		return
	}

	if err = global.GVA_DB.Delete(&[]bot.BotBanContent{}, "id in ?", ids).Error; err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("ids", IDs))
		return
	}

	for _, object := range objects {
		if deleteErr := cache.NewBotBanContentListCache(object.BotID).Release(); deleteErr != nil {
			global.GVA_LOG.Error("BotBanContentService", zap.Any("BotID", object.BotID))
		}
	}
	return err
}

// UpdateBotBanContent 更新机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (svc *BotBanContentService) UpdateBotBanContent(ctx context.Context, botBanContent bot.BotBanContent) (err error) {
	if err = global.GVA_DB.Model(&bot.BotBanContent{}).Where("id = ?", botBanContent.ID).Updates(&botBanContent).Error; err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("id", botBanContent.BotID))
		return
	}
	if deleteErr := cache.NewBotBanContentListCache(botBanContent.BotID).Release(); deleteErr != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("BotID", botBanContent.BotID))
	}
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
	db = db.Order("created_at desc")
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
	var botList []int64
	for _, botBanContent := range botBanContents {
		botList = append(botList, botBanContent.BotID)
	}

	var botMapper map[int64]bot.Bot
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_DB, botList); err != nil {
		return
	}

	for _, botBanContent := range botBanContents {
		botBanContent.BotName = botMapper[botBanContent.BotID].Name
	}
	return botBanContents, total, err
}
func (svc *BotBanContentService) GetBotBanContentPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
