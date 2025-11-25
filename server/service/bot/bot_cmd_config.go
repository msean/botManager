package bot

import (
	"context"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils"
	"go.uber.org/zap"
)

type BotCmdConfigService struct{}

// CreateBotCmdConfig 创建机器人命令配置记录
// Author [yourname](https://github.com/yourname)
func (botCmdConfigService *BotCmdConfigService) CreateBotCmdConfig(ctx context.Context, botCmdConfig *bot.BotCmdConfig) (err error) {
	err = global.GVA_DB.Create(botCmdConfig).Error
	return err
}

// DeleteBotCmdConfig 删除机器人命令配置记录
// Author [yourname](https://github.com/yourname)
func (botCmdConfigService *BotCmdConfigService) DeleteBotCmdConfig(ctx context.Context, ID string) (err error) {
	var id int
	if id, err = strconv.Atoi(ID); err != nil {
		return
	}
	var object bot.BotCmdConfig
	var has bool
	if has, err = utils.Get(global.GVA_DB, &object, utils.IDCond(id)); !has || err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("id", id), zap.Error(err))
		return
	}

	if err = global.GVA_DB.Delete(&bot.BotCmdConfig{}, "id = ?", ID).Error; err != nil {
		return
	}
	if deleteErr := cache.NewBotCmdCacheList(int(object.BotID)).Release(); deleteErr != nil {
		global.GVA_LOG.Error("botCmdConfigService", zap.Any("BotID", object.BotID))
	}
	if deleteErr := cache.NewBotCmdCache(object.BotID, object.Cmd).Release(); deleteErr != nil {
		global.GVA_LOG.Error("botCmdConfigService", zap.Any("BotID", object.BotID), zap.String("cmd", object.Cmd))
	}
	return err
}

// DeleteBotCmdConfigByIds 批量删除机器人命令配置记录
// Author [yourname](https://github.com/yourname)
func (botCmdConfigService *BotCmdConfigService) DeleteBotCmdConfigByIds(ctx context.Context, IDs []string) (err error) {
	ids := utils.StringsToIntsIgnoreError(IDs)
	var objects []bot.BotCmdConfig
	if err = utils.Find(global.GVA_DB, &objects, utils.NewInCond("id", utils.IntSliceToAnySlice(ids))); err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("ids", IDs), zap.Error(err))
		return
	}

	if err = global.GVA_DB.Delete(&[]bot.BotCmdConfig{}, "id in ?", IDs).Error; err != nil {
		return
	}

	for _, object := range objects {
		if deleteErr := cache.NewBotCmdCacheList(int(object.BotID)).Release(); deleteErr != nil {
			global.GVA_LOG.Error("BotBanContentService", zap.Any("BotID", object.BotID))
		}
		if deleteErr := cache.NewBotCmdCache(object.BotID, object.Cmd).Release(); deleteErr != nil {
			global.GVA_LOG.Error("botCmdConfigService", zap.Any("BotID", object.BotID), zap.String("cmd", object.Cmd))
		}
	}
	return err
}

// UpdateBotCmdConfig 更新机器人命令配置记录
// Author [yourname](https://github.com/yourname)
func (botCmdConfigService *BotCmdConfigService) UpdateBotCmdConfig(ctx context.Context, botCmdConfig bot.BotCmdConfig) (err error) {
	var object bot.BotCmdConfig
	var has bool
	if has, err = utils.Get(global.GVA_DB, &object, utils.IDCond(botCmdConfig.ID)); !has || err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("id", botCmdConfig.ID), zap.Error(err))
		return
	}
	spew.Dump(botCmdConfig)
	if err = global.GVA_DB.Model(&bot.BotCmdConfig{}).Where("id = ?", botCmdConfig.ID).Updates(&botCmdConfig).Error; err != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("id", botCmdConfig.BotID))
		return
	}
	if deleteErr := cache.NewBotCmdCacheList(int(object.BotID)).Release(); deleteErr != nil {
		global.GVA_LOG.Error("BotBanContentService", zap.Any("BotID", object.BotID))
	}
	if deleteErr := cache.NewBotCmdCache(object.BotID, object.Cmd).Release(); deleteErr != nil {
		global.GVA_LOG.Error("botCmdConfigService", zap.Any("BotID", object.BotID), zap.String("cmd", object.Cmd))
	}
	return err
}

// GetBotCmdConfig 根据ID获取机器人命令配置记录
// Author [yourname](https://github.com/yourname)
func (botCmdConfigService *BotCmdConfigService) GetBotCmdConfig(ctx context.Context, ID string) (botCmdConfig bot.BotCmdConfig, err error) {
	if err = global.GVA_DB.Where("id = ?", ID).First(&botCmdConfig).Error; err != nil {
		return
	}
	botModel, _, _ := dao.BotDao.FromBotID(global.GVA_DB, int(botCmdConfig.BotID))
	botCmdConfig.BotName = botModel.Name
	return
}

// GetBotCmdConfigInfoList 分页获取机器人命令配置记录
// Author [yourname](https://github.com/yourname)
func (botCmdConfigService *BotCmdConfigService) GetBotCmdConfigInfoList(ctx context.Context, info botReq.BotCmdConfigSearch) (list []*bot.BotCmdConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&bot.BotCmdConfig{})
	var botCmdConfigs []*bot.BotCmdConfig
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

	if err = db.Find(&botCmdConfigs).Error; err != nil {
		return
	}

	var botList []int
	for _, botBanContent := range botCmdConfigs {
		botList = append(botList, int(botBanContent.BotID))
	}

	var botMapper map[int]bot.Bot
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_DB, botList); err != nil {
		return
	}

	for _, botBanContent := range botCmdConfigs {
		botBanContent.BotName = botMapper[int(botBanContent.BotID)].Name
	}
	return botCmdConfigs, total, err
}

func (botCmdConfigService *BotCmdConfigService) GetBotCmdConfigPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
