package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/model/system/request"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

type BotMsgMassService struct{}

// CreateBotMsgMass 创建机器人群发记录
// Author [yourname](https://github.com/yourname)
// func (botMsgMassService *BotMsgMassService) CreateBotMsgMass(ctx context.Context, botMsgMass *bot.BotMsgMass) (err error) {
// 	err = global.GVA_MYSQL.Create(botMsgMass).Error
// 	return err
// }

// DeleteBotMsgMass 删除机器人群发记录
// Author [yourname](https://github.com/yourname)
// func (botMsgMassService *BotMsgMassService) DeleteBotMsgMass(ctx context.Context, ID string) (err error) {
// 	err = global.GVA_MYSQL.Delete(&bot.BotChatGroup{}, "id = ?", ID).Error
// 	return err
// }

// DeleteBotMsgMassByIds 批量删除机器人群发记录
// Author [yourname](https://github.com/yourname)
// func (botMsgMassService *BotMsgMassService) DeleteBotMsgMassByIds(ctx context.Context, IDs []string) (err error) {
// 	err = global.GVA_MYSQL.Delete(&[]bot.BotMsgMass{}, "id in ?", IDs).Error
// 	return err
// }

// UpdateBotMsgMass 更新机器人群发记录
// Author [yourname](https://github.com/yourname)
func (botMsgMassService *BotMsgMassService) UpdateBotMsgMass(ctx context.Context, botMsgMass bot.BotChatGroup) (err error) {
	err = global.GVA_MYSQL.Model(&bot.BotChatGroup{}).Where("id = ?", botMsgMass.ID).Updates(&botMsgMass).Error
	return err
}

// GetBotMsgMass 根据ID获取机器人群发记录
// Author [yourname](https://github.com/yourname)
func (botMsgMassService *BotMsgMassService) GetBotMsgMass(ctx context.Context, ID string) (botMsgMass bot.BotChatGroup, err error) {
	if err = global.GVA_MYSQL.Where("id = ?", ID).First(&botMsgMass).Error; err != nil {
		return
	}

	chatGroupModel, _, _ := dao.BotChatGroupDao.FromID(global.GVA_MYSQL, int(botMsgMass.ChatGroupID))
	botModel, _, _ := dao.BotDao.FromBotID(global.GVA_MYSQL, int(botMsgMass.BotID))

	botMsgMass.BotName = botModel.Name
	botMsgMass.ChatGroupName = chatGroupModel.ChatGroupName
	return
}

type BotGroupPair struct {
	BotID       int64
	ChatGroupID int64
}

func GetBelongs(ctx *gin.Context) (pairs []BotGroupPair, err error) {
	claimsVal := ctx.Value("claims")
	if claimsVal == nil {
		err = errors.New("获取用户失败")
		return
	}

	claims, ok := claimsVal.(*request.CustomClaims)
	if !ok {
		err = errors.New("解析claims失败")
		return
	}

	userID := claims.BaseClaims.ID
	authorityId := claims.AuthorityId

	isAdmin := authorityId == 220 || authorityId == 1

	if isAdmin {
		var list []BotGroupPair
		err = global.GVA_MYSQL.Model(&bot.BotChatGroup{}).
			Select("bot_id, chat_group_id").
			Scan(&list).Error
		return list, err
	}

	var classifies []bot.BotChatGroupClassify

	if err = global.GVA_MYSQL.Model(&bot.BotChatGroupClassify{}).
		Where(
			"users = ? OR users LIKE ? OR users LIKE ? OR users LIKE ?",
			fmt.Sprintf("%d", userID),
			fmt.Sprintf("%d,", userID)+"%",
			"%,"+fmt.Sprintf("%d,", userID)+"%",
			"%,"+fmt.Sprintf("%d", userID),
		).
		Find(&classifies).Error; err != nil {
		return
	}

	pairMap := make(map[string]BotGroupPair)

	for _, c := range classifies {
		if c.ChatGroups == "" {
			continue
		}

		items := strings.Split(c.ChatGroups, ",")

		for _, item := range items {
			parts := strings.Split(item, "_")
			if len(parts) != 2 {
				continue
			}

			botID, err1 := strconv.ParseInt(parts[0], 10, 64)
			groupID, err2 := strconv.ParseInt(parts[1], 10, 64)
			if err1 != nil || err2 != nil {
				continue
			}

			key := item
			pairMap[key] = BotGroupPair{
				BotID:       botID,
				ChatGroupID: groupID,
			}
		}
	}

	for _, v := range pairMap {
		pairs = append(pairs, v)
	}

	return
}

// GetBotMsgMassInfoList 分页获取机器人群发记录
// Author [yourname](https://github.com/yourname)
// GetBotMsgMassInfoList 分页获取机器人群发记录（带权限控制）
func (svc *BotMsgMassService) GetBotMsgMassInfoList(ctx *gin.Context, info botReq.BotMsgMassSearch) (list []*bot.BotChatGroup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	pairs, err := GetBelongs(ctx)
	if err != nil {
		return
	}

	if len(pairs) == 0 {
		return []*bot.BotChatGroup{}, 0, nil
	}

	db := global.GVA_MYSQL.Model(&bot.BotChatGroup{})

	var conditions []string
	var values []interface{}

	for _, p := range pairs {
		conditions = append(conditions, "(bot_id = ? AND chat_group_id = ?)")
		values = append(values, p.BotID, p.ChatGroupID)
	}

	db = db.Where(strings.Join(conditions, " OR "), values...)

	if info.BotID != nil {
		db = db.Where("bot_id = ?", *info.BotID)
	}
	if info.ChatGroupID != nil {
		db = db.Where("chat_group_id = ?", *info.ChatGroupID)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	var records []*bot.BotChatGroup

	if err = db.Count(&total).Error; err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Order("created_at desc").Find(&records).Error; err != nil {
		return
	}

	var botIDs []int64
	var groupIDs []int64

	for _, r := range records {
		botIDs = append(botIDs, r.BotID)
		groupIDs = append(groupIDs, r.ChatGroupID)
	}

	botMapper, err := dao.BotDao.MappByIDList(global.GVA_MYSQL, botIDs)
	if err != nil {
		return
	}

	chatGroupMapper, err := dao.BotChatGroupDao.MappByChatGroupIDList(global.GVA_MYSQL, groupIDs)
	if err != nil {
		return
	}

	for _, r := range records {
		if b, ok := botMapper[r.BotID]; ok {
			r.BotName = b.Name
		}
		if g, ok := chatGroupMapper[r.ChatGroupID]; ok {
			r.ChatGroupName = g.ChatGroupName
		}
	}

	return records, total, nil
}

func (botMsgMassService *BotMsgMassService) GetBotMsgMassPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// GetBotMassMsgRecordInfoList 分页获取群发历史记录记录
// Author [yourname](https://github.com/yourname)
// GetBotMassMsgRecordInfoList 分页获取群发历史记录记录
func (svc *BotMsgMassService) GetHistory(ctx *gin.Context, info botReq.BotMassMsgRecordSearch) (list []*bot.BotMassMsgRecord, total int64, err error) {

	// ✅ 改这里：获取组合权限
	pairs, err := GetBelongs(ctx)
	if err != nil {
		return
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_MYSQL.Model(&bot.BotMassMsgRecord{})

	if len(pairs) > 0 {
		var conditions []string
		var values []interface{}

		for _, p := range pairs {
			conditions = append(conditions, "(bot_id = ? AND chat_group_id = ?)")
			values = append(values, p.BotID, p.ChatGroupID)
		}

		db = db.Where(strings.Join(conditions, " OR "), values...)
	} else {
		return []*bot.BotMassMsgRecord{}, 0, nil
	}

	var records []*bot.BotMassMsgRecord

	// 条件过滤（保持不动）
	if info.BotID != nil {
		db = db.Where("bot_id = ?", *info.BotID)
	}
	if info.ChatGroupID != nil {
		db = db.Where("chat_group_id = ?", *info.ChatGroupID)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Order("created_at").Find(&records).Error; err != nil {
		return
	}

	// 获取机器人和群聊映射（保持不动）
	var botIDs, chatGroupIDs []int64
	for _, r := range records {
		botIDs = append(botIDs, r.BotID)
		chatGroupIDs = append(chatGroupIDs, r.ChatGroupID)
	}

	botMapper, err := dao.BotDao.MappByIDList(global.GVA_MYSQL, botIDs)
	if err != nil {
		return
	}
	chatGroupMapper, err := dao.BotChatGroupDao.MappByChatGroupIDList(global.GVA_MYSQL, chatGroupIDs)
	if err != nil {
		return
	}

	// 填充名称（保持不动）
	for _, r := range records {
		if b, ok := botMapper[r.BotID]; ok {
			r.BotName = b.Name
		}
		if g, ok := chatGroupMapper[r.ChatGroupID]; ok {
			r.ChatGroupName = g.ChatGroupName
		}
	}

	return records, total, nil
}

func (botMsgMassService *BotMsgMassService) SendBotMsgMass(
	ctx context.Context,
	req botReq.BotMsgMassSend,
) (err error) {

	var list []bot.BotChatGroup
	// var _ bot.Bot
	// var list []bot.BotChatGroup
	if err = global.GVA_MYSQL.
		Table("bot_chat_group AS g").
		Joins("LEFT JOIN bot b ON g.bot_id = b.bot_id").
		Where("g.id IN ? AND b.is_for_msg_mass = ?", req.IDs, 1). // ✅ 核心限制
		Select("g.*").
		Find(&list).Error; err != nil {
		return
	}

	// err = global.GVA_MYSQL.
	// 	Where("id IN ?", req.IDs).
	// 	Find(&list).Error
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return errors.New("未找到群发记录")
	}

	go func() {
		records := make([]bot.BotMassMsgRecord, 0, len(list))

		for _, item := range list {
			// 防止429
			time.Sleep(100 * time.Millisecond)
			// 👇 拼接 @成员
			remark := "发送成功"
			finalMsg := req.Msg
			if req.AtUsers {
				finalMsg = req.Msg + FormatMentionMembers(item.Members)
			}

			err := SendMassMsg(item.ChatGroupID, item.BotID, finalMsg)
			if err != nil {
				global.GVA_LOG.Error("SendBotMsgMass SendMassMsg", zap.Any("item", item), zap.Error(err))
				remark = err.Error()
			}
			records = append(records, bot.BotMassMsgRecord{
				BotID:       item.BotID,
				ChatGroupID: item.ChatGroupID,
				Msg:         req.Msg,
				Members:     item.Members,
				Remark:      remark,
			})
		}

		if err = global.GVA_MYSQL.Create(&records).Error; err != nil {
			global.GVA_LOG.Error("SendBotMsgMass Create", zap.Any("records", records), zap.Error(err))
		}
	}()
	return
}

func SendMassMsg(chatID int64, botID int64, content string) error {
	botModel, has, err := dao.BotDao.FromBotID(global.GVA_MYSQL, int(botID))
	if !has || err != nil {
		if !has {
			return fmt.Errorf("bot %d not found", botID)
		}
		return err
	}

	return bot_handler.HandleTexWithMarup(chatID, botModel.Token, content, nil)

}

func FormatMentionMembers(members string) string {
	if members == "" {
		return ""
	}

	arr := strings.Split(members, ",")
	res := make([]string, 0, len(arr))

	for _, m := range arr {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		if !strings.HasPrefix(m, "@") {
			m = "@" + m
		}
		res = append(res, m)
	}

	if len(res) == 0 {
		return ""
	}

	// 前面换两行，更好看
	return "\n\n" + strings.Join(res, " ")
}
