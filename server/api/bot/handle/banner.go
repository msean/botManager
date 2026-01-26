package handle

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func BanChatGroupContent(botModel bot.Bot, tgMsg botapi.Update) (find bool, err error) {
	if tgMsg.Message == nil {
		return
	}

	cacheObjects := cache.NewBotBanContentListCache(botModel.BotID)
	_, err = cache.CacheGetItem(cacheObjects)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int64("botID", botModel.BotID), zap.Error(err))
		return
	}

	messageText := tgMsg.Message.Text

	for _, rule := range cacheObjects.Objects {
		if strings.Contains(messageText, rule.BanContent) {
			global.GVA_LOG.Info("found banned word",
				zap.String("word", rule.BanContent),
				zap.String("user", tgMsg.Message.From.UserName),
				zap.String("msg", messageText),
			)

			BanUser(botModel, tgMsg, constant.BanTypeWord)

			find = true
			return
		}
	}
	return
}

func BanChatGroupMem(botModel bot.Bot, tgMsg botapi.Update) (found bool, err error) {
	if tgMsg.Message == nil {
		return
	}

	chatID := tgMsg.Message.Chat.ID
	user := tgMsg.Message.From

	cacheObjects := cache.NewBotChatGroupBanMemCListCache(botModel.BotID, chatID)
	_, err = cache.CacheGetItem(cacheObjects)
	global.GVA_LOG.Info("BanChatGroupMem",
		zap.Any("cacheObjects", cacheObjects),
		zap.Error(err),
	)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int64("botID", botModel.BotID), zap.Error(err))
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

			BanUser(botModel, tgMsg, constant.BanTypeMem)

			found = true
			return
		}
	}
	return
}
func BanUser(botModel bot.Bot, tgMsg botapi.Update, _type int) (err error) {
	ctx := context.Background()
	userKey := fmt.Sprintf("ban_user:%d_%d_%d", botModel.BotID, tgMsg.Message.From.ID, tgMsg.Message.Chat.ID)
	lockKey := fmt.Sprintf("ban_lock:%d_%d_%d", botModel.BotID, tgMsg.Message.From.ID, tgMsg.Message.Chat.ID)

	lockTTL := 10 * time.Second
	ok, err := global.GVA_REDIS.SetNX(ctx, lockKey, 1, lockTTL).Result()
	if err != nil {
		global.GVA_LOG.Error("redis lock error", zap.Error(err))
		return err
	}
	if !ok {
		global.GVA_LOG.Info("user is being processed, skip", zap.Int64("botID", botModel.BotID), zap.Int64("userID", tgMsg.Message.From.ID))
		return nil
	}
	defer global.GVA_REDIS.Del(ctx, lockKey) // 确保释放锁

	// 30秒内重复处理判断
	exists, _ := global.GVA_REDIS.Exists(ctx, userKey).Result()
	if exists > 0 {
		global.GVA_LOG.Info("user recently banned, skip", zap.Int64("botID", botModel.BotID), zap.Int64("userID", tgMsg.Message.From.ID))
		return nil
	}

	var sysCnf *cache.SysCnfCache
	if sysCnf, err = cache.LoadSyscnf(constant.SysCnfUserBanDuritonKey, true, constant.DefaultUserBanDuriton); err != nil {
		return
	}
	durationMinutes, _ := strconv.Atoi(sysCnf.Value)

	var botHandler *bot_handler.Bot
	if botHandler, err = bot_handler.NewBot(botModel.Token); err != nil {
		global.GVA_LOG.Error("BanUser NewBot", zap.Any("bot", botModel), zap.Error(err))
		return
	}

	chatID := tgMsg.Message.Chat.ID
	messageID := tgMsg.Message.MessageID
	var banErr error
	until := time.Now().Add(time.Duration(durationMinutes) * time.Minute).Unix()

	if banErr = botHandler.BanUser(chatID, tgMsg.Message.From.ID, until); banErr != nil {
		global.GVA_LOG.Error("ban user failed", zap.Error(banErr))
	} else {
		global.GVA_LOG.Info("ban user success",
			zap.Int64("chatID", chatID),
			zap.Int64("user_id", tgMsg.Message.From.ID),
			zap.Int64("until", until),
		)
	}

	if chatID != 0 && messageID != 0 && banErr == nil {
		if deleteErr := botHandler.DeleteMsg(chatID, messageID); deleteErr != nil {
			global.GVA_LOG.Error("delete msg error",
				zap.Int64("chatID", chatID),
				zap.Int64("user_id", tgMsg.Message.From.ID),
				zap.Int64("until", until),
				zap.Error(deleteErr),
			)
		}
	}

	remark := ""
	if banErr != nil {
		remark = banErr.Error()
	}

	_liftingTime := time.Now().Add(time.Duration(durationMinutes) * time.Minute)
	record := bot.BanRecord{
		BotID:       botModel.BotID,
		UserID:      tgMsg.Message.From.ID,
		UserName:    tgMsg.Message.From.UserName,
		ChatID:      chatID,
		ChatName:    tgMsg.Message.Chat.Title,
		BanDuration: int64(durationMinutes),
		Remark:      remark,
		BanType:     _type,
		FullName:    fmt.Sprintf("%s%s", tgMsg.Message.From.FirstName, tgMsg.Message.From.LastName),
		Msg:         tgMsg.Message.Text,
		LiftingTime: &_liftingTime,
		Status:      1, // 封禁中
	}
	if err := global.GVA_MYSQL.Create(&record).Error; err != nil {
		global.GVA_LOG.Error("failed to insert BanRecord", zap.Any("record", record), zap.Error(err))
		return err
	}
	global.GVA_REDIS.Set(ctx, userKey, 1, 30*time.Second)

	return nil
}
