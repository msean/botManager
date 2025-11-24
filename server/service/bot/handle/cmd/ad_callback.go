package cmd

import (
	"context"
	"encoding/json"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
)

func HandleAdCancel(chatID int64, userID int64, updateID int, token string, botID int64, msgID int) error {
	ctx := context.Background()
	bot, _ := tgbotapi.NewBotAPI(token)

	// 1. 先判断是否已经创建订单（最关键！！！）
	var count int64
	global.GVA_DB.
		Model(&recharge.UserRechargeRecord{}).
		Where("bot_id = ? AND user_id = ? AND update_id = ?", botID, userID, updateID).
		Count(&count)

	if count > 0 {
		bot.Send(tgbotapi.NewMessage(chatID,
			"⚠️ 订单已经创建，如需更改发布，请重新下单。"))
		del := tgbotapi.NewDeleteMessage(chatID, msgID)
		bot.Send(del)
		return nil
	}

	// 2. 判断草稿是否存在
	draftKey := cache.AdDraftCacheKey(botID, userID, int64(updateID))
	val, _ := global.GVA_REDIS.Get(ctx, draftKey).Result()

	if val == "" {
		// 草稿不存在 = 超时
		bot.Send(tgbotapi.NewMessage(chatID,
			"⏱️ 发布请求已超时，请重新提交内容。"))
		return nil
	}

	// 3. 正常取消
	global.GVA_REDIS.Del(ctx, draftKey)

	del := tgbotapi.NewDeleteMessage(chatID, msgID)
	bot.Send(del)

	bot.Send(tgbotapi.NewMessage(chatID, "❌ 已取消发布。"))

	return nil
}

func HandleAdConfirm(chatID int64, userID int64, updateID int64, token string, botID int64, msgID int) error {
	ctx := context.Background()

	draftKey := cache.AdDraftCacheKey(botID, userID, int64(updateID))

	val, err := global.GVA_REDIS.Get(ctx, draftKey).Result()
	if err != nil || val == "" {
		del := tgbotapi.NewDeleteMessage(chatID, msgID)
		bot, _ := tgbotapi.NewBotAPI(token)
		bot.Send(del)
		bot_handler.SendTextMessage(chatID, token, "❌ 此发布请求已过期，请重新发送内容。")
		return nil
	}

	// 写入订单
	rec := recharge.UserRechargeRecord{
		BotID:           botID,
		PublishTimes:    1,
		StartTime:       time.Now(),
		PublishInterval: 30,
		PublishContent:  val,
		Status:          1, // 创建
		UserID:          userID,
		UpdateID:        updateID,
	}
	global.GVA_DB.Create(&rec)

	global.GVA_REDIS.Del(ctx, draftKey)

	bot_handler.SendTextMessage(chatID, token, "✅ 广告订单创建成功，请前往后台完成支付。")

	// 查询所有频道
	var channels []bot.BotChannel
	global.GVA_DB.Where("bot_id = ?", botID).Find(&channels)

	// 群发
	var medias []MediaItem
	json.Unmarshal([]byte(val), &medias)
	for _, ch := range channels {
		AdSend(token, ch.ChannelID, medias, nil)
	}
	return nil
}
