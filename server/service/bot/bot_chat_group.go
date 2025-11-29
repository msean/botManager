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

type BotChatGroupService struct{}

// CreateBotChatGroup 创建机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) CreateBotChatGroup(ctx context.Context, botChatGroup *bot.BotChatGroup) (err error) {
	err = global.GVA_DB.Create(botChatGroup).Error
	return err
}

// DeleteBotChatGroup 删除机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) DeleteBotChatGroup(ctx context.Context, ID string) (err error) {
	var id int
	if id, err = strconv.Atoi(ID); err != nil {
		return
	}
	var object bot.BotChatGroup
	var has bool
	if has, err = utils.Get(global.GVA_DB, &object, utils.IDCond(ID)); !has || err != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("id", id), zap.Error(err))
		return
	}
	if err = global.GVA_DB.Delete(&bot.BotChatGroup{}, "id = ?", id).Error; err != nil {
		return
	}
	if deleteErr := cache.NewBotChatGroupCache(object).Release(); deleteErr != nil {
		global.GVA_LOG.Error("ReleaseBotChatGroup", zap.Any("BotID", object.BotID), zap.Int64("ChatGroupID", object.ChatGroupID))
	}
	return err
}

// DeleteBotChatGroupByIds 批量删除机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) DeleteBotChatGroupByIds(ctx context.Context, IDs []string) (err error) {
	ids := utils.StringsToIntsIgnoreError(IDs)
	var objects []bot.BotChatGroup
	if err = utils.Find(global.GVA_DB, &objects, utils.NewInCond("id", utils.IntSliceToAnySlice(ids))); err != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("ids", IDs), zap.Error(err))
		return
	}

	if err = global.GVA_DB.Delete(&[]bot.BotChatGroup{}, "id in ?", IDs).Error; err != nil {
		return
	}
	for _, object := range objects {
		if deleteErr := cache.NewBotChatGroupCache(object).Release(); deleteErr != nil {
			global.GVA_LOG.Error("ReleaseBotChatGroup", zap.Any("BotID", object.BotID))
		}
	}
	return err
}

// UpdateBotChatGroup 更新机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) UpdateBotChatGroup(ctx context.Context, botChatGroup bot.BotChatGroup) (err error) {
	err = global.GVA_DB.Model(&bot.BotChatGroup{}).Where("id = ?", botChatGroup.ID).Updates(&botChatGroup).Error
	return err
}

// GetBotChatGroup 根据ID获取机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) GetBotChatGroup(ctx context.Context, ID string) (botChatGroup bot.BotChatGroup, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&botChatGroup).Error
	return
}

// GetBotChatGroupInfoList 分页获取机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) GetBotChatGroupInfoList(ctx context.Context, info botReq.BotChatGroupSearch) (list []*bot.BotChatGroup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&bot.BotChatGroup{})
	var botChatGroups []*bot.BotChatGroup
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
	db = db.Order("created_at desc")

	if err = db.Find(&botChatGroups).Error; err != nil {
		return
	}

	var botList []int64
	for _, object := range botChatGroups {
		botList = append(botList, object.BotID)
	}

	var botMapper map[int64]bot.Bot
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_DB, botList); err != nil {
		return
	}
	for _, object := range botChatGroups {
		object.BotName = botMapper[object.BotID].Name
	}

	return botChatGroups, total, err
}
func (botChatGroupService *BotChatGroupService) GetBotChatGroupPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
