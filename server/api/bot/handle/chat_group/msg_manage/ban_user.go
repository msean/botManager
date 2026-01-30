package msgmanage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func BanUser(botModel bot.Bot, tgMsg botapi.Update, _type int, duration time.Duration) (err error) {
	var durationMinutes int
	if duration > 0 {
		durationMinutes = int(duration.Minutes())
	} else {
		var sysCnf *cache.SysCnfCache
		if sysCnf, err = cache.LoadSyscnf(constant.SysCnfUserBanDuritonKey, true, constant.DefaultUserBanDuriton); err != nil {
			return
		}
		durationMinutes, _ = strconv.Atoi(sysCnf.Value)
	}
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
