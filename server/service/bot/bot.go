package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/utils/bot_handler.go"
	"go.uber.org/zap"
)

type BotService struct{}

// CreateBot 创建机器人记录
// Author [yourname](https://github.com/yourname)
func (svc *BotService) CreateBot(ctx context.Context, botModel *bot.Bot) (err error) {
	parts := strings.Split(botModel.Token, ":")
	if len(parts) < 2 {
		return errors.New("输入的token不合法")
	}

	botID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return errors.New("输入的token不合法")
	}
	botModel.BotID = int(botID)

	// 检查是否存在
	if _, exist, err := dao.BotDao.FromBotID(global.GVA_DB, int(botID)); err != nil {
		return err
	} else if exist {
		return errors.New("已经存在相同的botID")
	}

	// 创建记录
	if err := global.GVA_DB.Create(botModel).Error; err != nil {
		return err
	}

	webhookURL := fmt.Sprintf("%s/bot/webhook/%s", global.GVA_CONFIG.System.RouterPrefix, botModel.Token)

	go func() {
		if err := bot_handler.RegisterWebhook(botModel.Token, webhookURL); err != nil {
			global.GVA_LOG.Error("register webhook failed", zap.String("url", webhookURL), zap.Error(err))
		} else {
			global.GVA_LOG.Info("register webhook success", zap.String("url", webhookURL))
		}
	}()

	return nil
}

// DeleteBot 删除机器人记录
// Author [yourname](https://github.com/yourname)
func (svc *BotService) DeleteBot(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&bot.Bot{}, "bot_id = ?", ID).Error
	return err
}

// DeleteBotByIds 批量删除机器人记录
// Author [yourname](https://github.com/yourname)
func (svc *BotService) DeleteBotByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]bot.Bot{}, "bot_id in ?", IDs).Error
	return err
}

// UpdateBot 更新机器人记录
// Author [yourname](https://github.com/yourname)
func (svc *BotService) UpdateBot(ctx context.Context, botModel bot.Bot) (err error) {
	err = global.GVA_DB.Model(&bot.Bot{}).Where("bot_id = ?", botModel.BotID).Updates(&botModel).Error
	return err
}

// GetBot 根据ID获取机器人记录
// Author [yourname](https://github.com/yourname)
func (svc *BotService) GetBot(ctx context.Context, ID string) (bot_mgr bot.Bot, err error) {
	err = global.GVA_DB.Where("bot_id = ?", ID).First(&bot_mgr).Error
	return
}

// GetBotInfoList 分页获取机器人记录
// Author [yourname](https://github.com/yourname)
func (svc *BotService) GetBotInfoList(ctx context.Context, info botReq.BotSearch) (list []bot.Bot, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&bot.Bot{})
	var bot_mgrs []bot.Bot
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Token != "" {
		db = db.Where("token LIKE ?", "%"+info.Token+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&bot_mgrs).Error
	return bot_mgrs, total, err
}
func (svc *BotService) GetBotPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
