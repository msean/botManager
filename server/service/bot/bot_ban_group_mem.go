package bot

import (
	"context"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
)

type BotBanGroupMemService struct{}

// CreateBotBanGroupMem 创建封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) CreateBotBanGroupMem(ctx context.Context, botBanGroupMem *bot.BotBanGroupMem) (err error) {
	err = global.GVA_DB.Create(botBanGroupMem).Error
	return err
}

// DeleteBotBanGroupMem 删除封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) DeleteBotBanGroupMem(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&bot.BotBanGroupMem{}, "id = ?", ID).Error
	return err
}

// DeleteBotBanGroupMemByIds 批量删除封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) DeleteBotBanGroupMemByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]bot.BotBanGroupMem{}, "id in ?", IDs).Error
	return err
}

// UpdateBotBanGroupMem 更新封禁成员设置记录
// Author [yourname](https://github.com/yourname)
func (botBanGroupMemService *BotBanGroupMemService) UpdateBotBanGroupMem(ctx context.Context, botBanGroupMem bot.BotBanGroupMem) (err error) {
	err = global.GVA_DB.Model(&bot.BotBanGroupMem{}).Where("id = ?", botBanGroupMem.ID).Updates(&botBanGroupMem).Error
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

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

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
