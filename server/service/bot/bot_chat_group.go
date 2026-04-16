package bot

import (
	"context"
	"strconv"
	"strings"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils"
)

type BotChatGroupService struct{}

// CreateBotChatGroup 创建机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) CreateBotChatGroup(ctx context.Context, botChatGroup *bot.BotChatGroup) (err error) {
	err = global.GVA_MYSQL.Create(botChatGroup).Error
	return err
}

// DeleteBotChatGroup 删除机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) DeleteBotChatGroup(ctx context.Context, ID string) (err error) {
	var id int
	if id, err = strconv.Atoi(ID); err != nil {
		return
	}
	if err = cache.ReleaseBotChatGroup(id); err != nil {
		return
	}
	if err = global.GVA_MYSQL.Delete(&bot.BotChatGroup{}, "id = ?", id).Error; err != nil {
		return
	}
	return err
}

// DeleteBotChatGroupByIds 批量删除机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) DeleteBotChatGroupByIds(ctx context.Context, IDs []string) (err error) {
	ids := utils.StringsToIntsIgnoreError(IDs)
	for _, id := range ids {
		cache.ReleaseBotChatGroup(id)
	}
	if err = global.GVA_MYSQL.Delete(&[]bot.BotChatGroup{}, "id in (?)", IDs).Error; err != nil {
		return
	}

	return err
}

// UpdateBotChatGroup 更新机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) UpdateBotChatGroup(ctx context.Context, botChatGroup bot.BotChatGroup) (err error) {
	if err = cache.ReleaseBotChatGroup(int(botChatGroup.ID)); err != nil {
		return
	}
	return global.GVA_MYSQL.Model(&bot.BotChatGroup{}).Where("id = ?", botChatGroup.ID).Updates(&botChatGroup).Error
}

// GetBotChatGroup 根据ID获取机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) GetBotChatGroup(ctx context.Context, ID string) (botChatGroup bot.BotChatGroup, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&botChatGroup).Error
	return
}

// GetBotChatGroupInfoList 分页获取机器人群组列表记录
// Author [yourname](https://github.com/yourname)
func (botChatGroupService *BotChatGroupService) GetBotChatGroupInfoList(ctx context.Context, info botReq.BotChatGroupSearch) (list []*bot.BotChatGroup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&bot.BotChatGroup{})
	var botChatGroups []*bot.BotChatGroup
	if info.BotID != 0 {
		db = db.Where("bot_id = ?", info.BotID)
	}
	if info.Name != "" {
		db = db.Where("chat_group_name LIKE ?", "%"+info.Name+"%")
	}
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
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList); err != nil {
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

func uniqueInt64(arr []int64) []int64 {
	set := make(map[int64]struct{})
	res := make([]int64, 0, len(arr))

	for _, v := range arr {
		if _, ok := set[v]; !ok {
			set[v] = struct{}{}
			res = append(res, v)
		}
	}
	return res
}

func (s *BotChatGroupService) ClassfyList(info botReq.BotChatGroupClassifySearch) (
	list []bot.BotChatGroupClassify,
	botMapper map[int64]string,
	chatGroupMapper map[int64]string,
	userMapper map[int64]string,
	total int64,
	err error,
) {

	db := global.GVA_MYSQL.Model(&bot.BotChatGroupClassify{})

	// 总数
	if err = db.Count(&total).Error; err != nil {
		return
	}

	offset := (info.Page - 1) * info.PageSize
	if err = db.Limit(info.PageSize).
		Offset(offset).
		Order("id desc").
		Find(&list).Error; err != nil {
		return
	}

	var botList []int64
	var chatGroupList []int64
	var userList []int64

	for _, item := range list {
		chatGroupIDList := strings.Split(item.ChatGroups, ",")
		for _, raw := range chatGroupIDList {
			if raw == "" {
				continue
			}

			parts := strings.Split(raw, "_")
			var groupID int64
			var botID int64
			if len(parts) == 2 {
				if gid, e := strconv.ParseInt(parts[1], 10, 64); e == nil {
					groupID = gid
				}
				if bid, e := strconv.ParseInt(parts[0], 10, 64); e == nil {
					botID = bid
				}
			} else {
				if gid, e := strconv.ParseInt(raw, 10, 64); e == nil {
					groupID = gid
				}
			}
			chatGroupList = append(chatGroupList, groupID)
			botList = append(botList, botID)

		}

		userIDList := strings.Split(item.Users, ",")

		for _, raw := range userIDList {
			if raw == "" {
				continue
			}
			if uid, e := strconv.ParseInt(raw, 10, 64); e == nil {
				userList = append(userList, uid)
			}
		}
	}

	chatGroupList = uniqueInt64(chatGroupList)
	userList = uniqueInt64(userList)

	if len(chatGroupList) > 0 {
		if chatGroupMapper, err = dao.BotChatGroupDao.
			MappNameByChatGroupIDList(global.GVA_MYSQL, chatGroupList); err != nil {
			return
		}
	} else {
		chatGroupMapper = map[int64]string{}
	}

	if len(userList) > 0 {
		if userMapper, err = dao.SysDao.
			NameMapperFromIDList(global.GVA_MYSQL, userList); err != nil {
			return
		}
	} else {
		userMapper = map[int64]string{}
	}

	if len(botList) > 0 {
		if botMapper, err = dao.BotDao.
			NameMappByIDList(global.GVA_MYSQL, botList); err != nil {
			return
		}
	} else {
		botMapper = map[int64]string{}

	}

	return
}

// 创建 / 更新
func (s *BotChatGroupService) SaveClassify(data bot.BotChatGroupClassify) error {
	db := global.GVA_MYSQL
	if data.ID == 0 {
		return db.Create(&data).Error
	}

	if data.Refresh {
		return db.Save(&data).Error
	}

	var old bot.BotChatGroupClassify
	if err := db.Where("id = ?", data.ID).First(&old).Error; err != nil {
		return err
	}

	groupMap := make(map[string]struct{})

	if old.ChatGroups != "" {
		for _, v := range strings.Split(old.ChatGroups, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				groupMap[v] = struct{}{}
			}
		}
	}

	if data.ChatGroups != "" {
		for _, v := range strings.Split(data.ChatGroups, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				groupMap[v] = struct{}{}
			}
		}
	}

	var result []string
	for k := range groupMap {
		result = append(result, k)
	}

	newChatGroups := strings.Join(result, ",")

	return db.Model(&old).Update("chat_groups", newChatGroups).Error
}

// 删除
func (s *BotChatGroupService) DeleteClassify(ids []uint) error {
	return global.GVA_MYSQL.Delete(&bot.BotChatGroupClassify{}, "id in (?)", ids).Error
}

// 列表
func (s *BotChatGroupService) ClassifyChoice() (list []bot.BotChatGroupClassify, err error) {
	db := global.GVA_MYSQL.Model(&bot.BotChatGroupClassify{})
	err = db.Order("id desc").Find(&list).Error
	return
}
