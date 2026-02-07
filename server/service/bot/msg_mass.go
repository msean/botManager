package bot

import (
	"context"
	"errors"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
)

type BotMsgMassService struct{}

// CreateBotMsgMass 创建机器人群发记录
// Author [yourname](https://github.com/yourname)
func (botMsgMassService *BotMsgMassService) CreateBotMsgMass(ctx context.Context, botMsgMass *bot.BotMsgMass) (err error) {
	err = global.GVA_MYSQL.Create(botMsgMass).Error
	return err
}

// DeleteBotMsgMass 删除机器人群发记录
// Author [yourname](https://github.com/yourname)
func (botMsgMassService *BotMsgMassService) DeleteBotMsgMass(ctx context.Context, ID string) (err error) {
	err = global.GVA_MYSQL.Delete(&bot.BotMsgMass{}, "id = ?", ID).Error
	return err
}

// DeleteBotMsgMassByIds 批量删除机器人群发记录
// Author [yourname](https://github.com/yourname)
func (botMsgMassService *BotMsgMassService) DeleteBotMsgMassByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]bot.BotMsgMass{}, "id in ?", IDs).Error
	return err
}

// UpdateBotMsgMass 更新机器人群发记录
// Author [yourname](https://github.com/yourname)
func (botMsgMassService *BotMsgMassService) UpdateBotMsgMass(ctx context.Context, botMsgMass bot.BotMsgMass) (err error) {
	err = global.GVA_MYSQL.Model(&bot.BotMsgMass{}).Where("id = ?", botMsgMass.ID).Updates(&botMsgMass).Error
	return err
}

// GetBotMsgMass 根据ID获取机器人群发记录
// Author [yourname](https://github.com/yourname)
func (botMsgMassService *BotMsgMassService) GetBotMsgMass(ctx context.Context, ID string) (botMsgMass bot.BotMsgMass, err error) {
	err = global.GVA_MYSQL.Where("id = ?", ID).First(&botMsgMass).Error
	return
}

// GetBotMsgMassInfoList 分页获取机器人群发记录
// Author [yourname](https://github.com/yourname)
func (botMsgMassService *BotMsgMassService) GetBotMsgMassInfoList(ctx context.Context, info botReq.BotMsgMassSearch) (list []*bot.BotMsgMass, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&bot.BotMsgMass{})
	var botMsgMasss []*bot.BotMsgMass
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if err = db.Find(&botMsgMasss).Error; err != nil {
		return
	}

	var botList []int64
	var chatGroupList []int64
	for _, object := range botMsgMasss {
		botList = append(botList, object.BotID)
		chatGroupList = append(chatGroupList, object.ChatGroupID)
	}

	var botMapper map[int64]bot.Bot
	var chatGroupMapper map[int64]bot.BotChatGroup
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList); err != nil {
		return
	}

	if chatGroupMapper, err = dao.BotChatGroupDao.MappByChatGroupIDList(global.GVA_MYSQL, chatGroupList); err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	for _, object := range botMsgMasss {
		object.BotName = botMapper[object.BotID].Name
		object.ChatGroupName = chatGroupMapper[object.ChatGroupID].ChatGroupName
	}
	return botMsgMasss, total, err
}

func (botMsgMassService *BotMsgMassService) GetBotMsgMassPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// UpdateBotMsgMass 更新机器人群发记录
// Author [yourname](https://github.com/yourname)
// SendBotMsgMass 群发（仅保存记录，不发 Telegram）
func (botMsgMassService *BotMsgMassService) SendBotMsgMass(
	ctx context.Context,
	req botReq.BotMsgMassSend,
) (err error) {

	var list []bot.BotMsgMass

	// 1️⃣ 查出被选中的群发配置
	err = global.GVA_MYSQL.
		Where("id IN ?", req.IDs).
		Find(&list).Error
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return errors.New("未找到群发记录")
	}

	// 2️⃣ 组装群发历史记录
	records := make([]bot.BotMassMsgRecord, 0, len(list))

	for _, item := range list {
		records = append(records, bot.BotMassMsgRecord{
			BotID:       item.BotID,
			ChatGroupID: item.ChatGroupID,
			Msg:         req.Msg,
			Members:     item.Members,
			Remark:      "",
		})
	}

	// 3️⃣ 批量写入历史表
	err = global.GVA_MYSQL.
		Create(&records).Error

	return err
}
