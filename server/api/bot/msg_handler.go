package bot

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/common/response"
	"github.com/msean/botmanager/server/utils/bot_handler.go"
	"go.uber.org/zap"
)

type BotMsgHandler struct{}

func (api *BotMsgHandler) Handle(c *gin.Context) {
	botIDStr := c.Param("botUUID")
	botID, err := strconv.Atoi(botIDStr)
	if err != nil {
		response.BotBadRequest(c, "invalid botID")
		return
	}

	body, _ := io.ReadAll(c.Request.Body)
	c.Status(200) // Telegram 要求 webhook 必须快速响应，否则会重试
	global.GVA_LOG.Info("receive telegram webhook", zap.String("body", string(body)))

	// 解析 Telegram 消息
	var tgMsg bot.TelegramMessage
	if err := json.Unmarshal(body, &tgMsg); err != nil {
		global.GVA_LOG.Error("invalid telegram message", zap.Error(err))
		return
	}
	if tgMsg.Message.Text == "" {
		return // 不是文本消息（可能是 sticker、照片 等）
	}

	// 查找机器人信息
	botModel, has, err := dao.BotDao.FromBotID(global.GVA_DB, botID)
	if err != nil || !has {
		global.GVA_LOG.Error("bot not found", zap.Error(err))
		return
	}

	banWords, err := dao.BotDao.ListBotBannerContentByID(global.GVA_DB, botModel.BotID)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Error(err))
		return
	}

	messageText := strings.ToLower(tgMsg.Message.Text)

	for _, rule := range banWords {
		if strings.Contains(messageText, strings.ToLower(rule.BanContent)) {
			global.GVA_LOG.Info("found banned word",
				zap.String("word", rule.BanContent),
				zap.String("user", tgMsg.Message.From.UserName),
			)

			// 获取封禁时长
			param, err := dao.SysParamsDao.FromKey(global.GVA_DB, global.UserBanDuritonKey, global.DefaultUserBanDuriton)
			if err != nil {
				global.GVA_LOG.Error("get ban duration failed", zap.Error(err))
				return
			}
			durationMinutes, _ := strconv.Atoi(param.Value)

			go func() {
				// 发送api 封禁用户
				botHandler := bot_handler.NewBot(botModel.Token)
				if err := botHandler.BanUser(tgMsg.Message.Chat.ID, tgMsg.Message.From.ID, time.Duration(durationMinutes)*time.Minute); err != nil {
					global.GVA_LOG.Error("ban user failed", zap.Error(err))
				}
			}()

			record := bot.BanRecord{
				BotID:       botModel.BotID,
				UserID:      tgMsg.Message.From.ID,
				UserName:    tgMsg.Message.From.UserName,
				ChatID:      tgMsg.Message.Chat.ID,
				ChatName:    tgMsg.Message.Chat.Title,
				BanDuration: int64(durationMinutes),
			}
			if err := global.GVA_DB.Create(&record).Error; err != nil {
				global.GVA_LOG.Error("failed to insert BanRecord", zap.Error(err))
			}

			break
		}
	}
}
