package handle

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

func BanChatGroupContent(botModel bot.Bot, tgMsg tgbotapi.Update) (find bool, err error) {
	if tgMsg.Message == nil {
		return
	}

	cacheObjects := cache.NewBotBanContentListCache(botModel.BotID)
	_, err = cache.CacheGetItem(cacheObjects)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int("botID", botModel.BotID), zap.Error(err))
		return
	}

	messageText := tgMsg.Message.Text
	global.GVA_LOG.Info("msg", zap.String("bot", botModel.Name), zap.String("msg", tgMsg.Message.Text))

	for _, rule := range cacheObjects.Objects {
		if strings.Contains(messageText, rule.BanContent) {
			global.GVA_LOG.Info("found banned word",
				zap.String("word", rule.BanContent),
				zap.String("user", tgMsg.Message.From.UserName),
				zap.String("msg", messageText),
			)

			BanUser(botModel, tgMsg, global.BanTypeWord)

			find = true
			return
		}
	}
	return
}

func BanChatGroupMem(botModel bot.Bot, tgMsg tgbotapi.Update) (found bool, err error) {
	if tgMsg.Message == nil {
		return
	}

	chatID := tgMsg.Message.Chat.ID
	user := tgMsg.Message.From

	cacheObjects := cache.NewBotChatGroupBanMemCListCache(botModel.BotID)
	_, err = cache.CacheGetItem(cacheObjects)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int("botID", botModel.BotID), zap.Error(err))
		return
	}
	fullName := fmt.Sprintf("%s%s", user.FirstName, user.LastName)

	for _, ban := range cacheObjects.Objects {
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

			BanUser(botModel, tgMsg, global.BanTypeMem)

			found = true
			return
		}
	}
	return
}

func BanUser(botModel bot.Bot, tgMsg tgbotapi.Update, _type int) (err error) {
	var sysCnf *cache.SysCnfCache
	if sysCnf, err = cache.LoadSyscnf(global.SysCnfUserBanDuritonKey, true, global.DefaultUserBanDuriton); err != nil {
		return
	}
	durationMinutes, _ := strconv.Atoi(sysCnf.Value)

	// 发送api 封禁用户
	chatID := tgMsg.Message.Chat.ID
	messageID := tgMsg.Message.MessageID

	var botHandler *bot_handler.Bot
	if botHandler, err = bot_handler.NewBot(botModel.Token); err != nil {
		global.GVA_LOG.Error("BanUser NewBot", zap.Any("bot", botModel), zap.Error(err))
		return
	}
	var banErr error
	until := time.Now().Add(time.Duration(durationMinutes) * time.Minute).Unix()
	if banErr = botHandler.BanUser(tgMsg.Message.Chat.ID, tgMsg.Message.From.ID, until); banErr != nil {
		global.GVA_LOG.Error("ban user failed", zap.Error(err))
	} else {
		global.GVA_LOG.Info("ban user success",
			zap.Int64("chatID", tgMsg.Message.Chat.ID),
			zap.Int64("user_id", tgMsg.Message.Chat.ID),
			zap.Int64("util", until),
		)
	}

	if chatID != 0 && messageID != 0 && banErr == nil {
		if deleteErr := botHandler.DeleteMsg(chatID, messageID); deleteErr != nil {
			global.GVA_LOG.Error("delete msg error",
				zap.Int64("chatID", tgMsg.Message.Chat.ID),
				zap.Int64("user_id", tgMsg.Message.Chat.ID),
				zap.Int64("util", until),
				zap.Error(deleteErr),
			)
		}
	}

	remark := ""
	if banErr != nil {
		remark = banErr.Error()
	}
	msg := ""
	if _type == global.BanTypeWord {
		msg = tgMsg.Message.Text
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
		Msg:         msg,
	}
	if err := global.GVA_DB.Create(&record).Error; err != nil {
		global.GVA_LOG.Error("failed to insert BanRecord", zap.Any("record", record), zap.Error(err))
	}
	return
}
