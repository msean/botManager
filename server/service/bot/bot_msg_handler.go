package bot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/system"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

type BotMsgHandlerSvc struct{}

func (svc *BotMsgHandlerSvc) Handle(c *gin.Context, botID int, body []byte) (err error) {

	global.GVA_LOG.Debug("receive telegram webhook", zap.String("body", string(body)))

	// 解析 Telegram 消息
	var tgMsg bot.TelegramMessage
	if err = json.Unmarshal(body, &tgMsg); err != nil {
		global.GVA_LOG.Error("invalid telegram message", zap.Error(err))
		return
	}
	if tgMsg.Message.Text == "" {
		return // 不是文本消息（可能是 sticker、照片 等）
	}

	// 查找机器人信息
	botModel, has, err := dao.BotDao.FromBotID(global.GVA_DB, botID)
	if err != nil || !has {
		global.GVA_LOG.Error("bot not found", zap.Int("botID", botModel.BotID), zap.Error(err))
		return
	}

	var find bool
	if find, err = svc.CheckBanContent(botModel, tgMsg); err != nil || find {
		return
	}

	if find, err = svc.CheckGroupMem(botModel, tgMsg); err != nil || find {
		return
	}

	return err
}

func (svc *BotMsgHandlerSvc) BanUser(botModel bot.Bot, tgMsg bot.TelegramMessage, durationMinutes int, _type int) (err error) {
	// 发送api 封禁用户
	botHandler := bot_handler.NewBot(botModel.Token)
	var banErr error
	until := time.Now().Add(time.Duration(durationMinutes) * time.Minute).Unix()
	if banErr = botHandler.BanUser(tgMsg.Message.Chat.ID, tgMsg.Message.From.ID, until); err != nil {
		global.GVA_LOG.Error("ban user failed", zap.Error(err))
	} else {
		global.GVA_LOG.Info("ban user success",
			zap.Int64("chatID", tgMsg.Message.Chat.ID),
			zap.Int64("user_id", tgMsg.Message.Chat.ID),
			zap.Int64("util", until),
		)
	}
	remark := ""
	if banErr != nil {
		remark = banErr.Error()
	}
	record := bot.BanRecord{
		BotID:       botModel.BotID,
		UserID:      tgMsg.Message.From.ID,
		UserName:    tgMsg.Message.From.UserName,
		ChatID:      tgMsg.Message.Chat.ID,
		ChatName:    tgMsg.Message.Chat.Title,
		BanDuration: int64(durationMinutes),
		Remark:      remark,
		BanType:     _type,
		FullName:    fmt.Sprintf("%s%s", tgMsg.Message.From.FirstName, tgMsg.Message.From.LastName),
	}
	if err := global.GVA_DB.Create(&record).Error; err != nil {
		global.GVA_LOG.Error("failed to insert BanRecord", zap.Any("record", record), zap.Error(err))
	}
	return
}

func (svc *BotMsgHandlerSvc) CheckBanContent(botModel bot.Bot, tgMsg bot.TelegramMessage) (find bool, err error) {
	var banContents []bot.BotBanContent
	banContents, err = dao.BotDao.ListBotBannerContentByID(global.GVA_DB, botModel.BotID)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int("botID", botModel.BotID), zap.Error(err))
		return
	}

	messageText := strings.ToLower(tgMsg.Message.Text)

	for _, rule := range banContents {
		if strings.Contains(messageText, strings.ToLower(rule.BanContent)) {
			global.GVA_LOG.Info("found banned word",
				zap.String("word", rule.BanContent),
				zap.String("user", tgMsg.Message.From.UserName),
			)

			// 获取封禁时长
			var param system.SysParams
			param, err = dao.SysParamsDao.FromKey(global.GVA_DB, global.UserBanDuritonKey, global.DefaultUserBanDuriton)
			if err != nil {
				global.GVA_LOG.Error("get ban duration failed", zap.Error(err))
				return
			}
			durationMinutes, _ := strconv.Atoi(param.Value)

			go svc.BanUser(botModel, tgMsg, durationMinutes, global.BanTypeWord)

			find = true
			return
		}
	}
	return
}

func (svc *BotMsgHandlerSvc) CheckGroupMem(botModel bot.Bot, tgMsg bot.TelegramMessage) (found bool, err error) {
	chatID := tgMsg.Message.Chat.ID
	user := tgMsg.Message.From

	banList, err := dao.BotGroupMemDao.ListByBotIDAndChatGroupID(global.GVA_DB, botModel.BotID, int(chatID))
	if err != nil {
		global.GVA_LOG.Error("fetch ban members failed", zap.Int("botID", botModel.BotID), zap.Int64("chatID", chatID), zap.Error(err))
		return
	}
	if len(banList) == 0 {
		return
	}

	fullName := fmt.Sprintf("%s%s", user.FirstName, user.LastName)

	for _, ban := range banList {
		banStr := strings.ToLower(strings.TrimSpace(ban.BanMemContent))
		if banStr == "" {
			continue
		}

		if strings.Contains(fullName, banStr) {
			global.GVA_LOG.Info("found banned member",
				zap.String("ban_word", ban.BanMemContent),
				zap.String("member", user.UserName),
				zap.Int64("chatID", chatID),
			)

			var param system.SysParams
			param, err = dao.SysParamsDao.FromKey(global.GVA_DB, global.UserBanDuritonKey, global.DefaultUserBanDuriton)
			if err != nil {
				global.GVA_LOG.Error("get ban duration failed", zap.Error(err))
				continue
			}
			durationMinutes, _ := strconv.Atoi(param.Value)

			go svc.BanUser(botModel, tgMsg, durationMinutes, global.BanTypeMem)

			found = true
			return
		}
	}

	return
}
