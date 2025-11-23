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

type BotBanGroupMemService struct{}

// CreateBotBanGroupMem 创建封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) CreateBotBanGroupMem(ctx context.Context, botBanGroupMem *bot.BotBanGroupMem) (err error) {
	if err = global.GVA_DB.Create(botBanGroupMem).Error; err != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("botBanGroupMem", botBanGroupMem), zap.Error(err))
		return
	}
	return err
}

// DeleteBotBanGroupMem 删除封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) DeleteBotBanGroupMem(ctx context.Context, ID string) (err error) {
	var id int
	if id, err = strconv.Atoi(ID); err != nil {
		return
	}
	var object bot.BotBanGroupMem
	var has bool
	if has, err = utils.Get(global.GVA_DB, &object, utils.IDCond(ID)); !has || err != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("id", id), zap.Error(err))
		return
	}
	if err = global.GVA_DB.Delete(&bot.BotBanGroupMem{}, "id = ?", ID).Error; err != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("id", id))
		return
	}

	if deleteErr := cache.NewBotChatGroupBanMemCListCache(int(object.BotID)).Release(); deleteErr != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("BotID", object.BotID), zap.Int64("ChatGroupID", object.ChatGroupID))
	}
	return err
}

// DeleteBotBanGroupMemByIds 批量删除封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) DeleteBotBanGroupMemByIds(ctx context.Context, IDs []string) (err error) {
	ids := utils.StringsToIntsIgnoreError(IDs)
	var objects []bot.BotBanGroupMem
	if err = utils.Find(global.GVA_DB, &objects, utils.NewInCond("id", utils.IntSliceToAnySlice(ids))); err != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("ids", IDs), zap.Error(err))
		return
	}

	if err = global.GVA_DB.Delete(&[]bot.BotBanGroupMem{}, "id in ?", ids).Error; err != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("ids", IDs))
		return
	}
	for _, object := range objects {
		if deleteErr := cache.NewBotChatGroupBanMemCListCache(int(object.BotID)).Release(); deleteErr != nil {
			global.GVA_LOG.Error("botBanGroupMemService", zap.Any("BotID", object.BotID))
		}
	}
	return err
}

// UpdateBotBanGroupMem 更新封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) UpdateBotBanGroupMem(ctx context.Context, botBanGroupMem bot.BotBanGroupMem) (err error) {
	if err = global.GVA_DB.Model(&bot.BotBanGroupMem{}).Where("id = ?", botBanGroupMem.ID).Updates(&botBanGroupMem).Error; err != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("id", botBanGroupMem.ID))
		return
	}
	if deleteErr := cache.NewBotChatGroupBanMemCListCache(int(botBanGroupMem.BotID)).Release(); deleteErr != nil {
		global.GVA_LOG.Error("botBanGroupMemService", zap.Any("botBanGroupMem", botBanGroupMem))
	}
	return err
}

// GetBotBanGroupMem 根据ID获取封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) GetBotBanGroupMem(ctx context.Context, ID string) (botBanGroupMem bot.BotBanGroupMem, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&botBanGroupMem).Error
	return
}

// GetBotBanGroupMemInfoList 分页获取封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) GetBotBanGroupMemInfoList(ctx context.Context, info botReq.BotBanGroupMemSearch) (list []*bot.BotBanGroupMem, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&bot.BotBanGroupMem{})
	var botBanGroupMems []*bot.BotBanGroupMem
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
	db = db.Order("created_at desc")

	if err = db.Find(&botBanGroupMems).Error; err != nil {
		return
	}

	var botList []int
	var chatGroupList []int
	for _, object := range botBanGroupMems {
		botList = append(botList, int(object.BotID))
		chatGroupList = append(chatGroupList, int(object.ChatGroupID))
	}

	var botMapper map[int]bot.Bot
	var chatGroupMapper map[int]bot.BotChatGroup
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_DB, botList); err != nil {
		return
	}

	if chatGroupMapper, err = dao.BotChatGroupDao.MappByChatGroupIDList(global.GVA_DB, chatGroupList); err != nil {
		return
	}

	for _, object := range botBanGroupMems {
		object.BotName = botMapper[int(object.BotID)].Name
		object.ChatGroupName = chatGroupMapper[int(object.ChatGroupID)].ChatGroupName
	}

	return botBanGroupMems, total, err
}
func (botBanGroupMemService *BotBanGroupMemService) GetBotBanGroupMemPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
