package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	rechargeSrv "github.com/msean/botmanager/server/service/recharge"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
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

func HandleAdConfirm(chatID int64, userID int64, updateID int64, token string, botID int64, msgID int, publishTimes int) (err error) {
	ctx := context.Background()

	draftKey := cache.AdDraftCacheKey(botID, userID, int64(updateID))
	var botHandler *bot_handler.Bot
	if botHandler, err = bot_handler.NewBot(token); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
		return
	}

	val, err := global.GVA_REDIS.Get(ctx, draftKey).Result()
	if err != nil || val == "" {

		if err = botHandler.DeleteMsg(chatID, msgID); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm DeleteMsg", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
		}
		if err = botHandler.SendTextMessage(chatID, token, "❌ 此发布请求已过期，请重新发送内容。"); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
		}
		return nil
	}

	var price float64
	pay := rechargeSrv.NewPay(botID, publishTimes)
	if price, err = pay.RandomPrice(); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("publishTimes", int64(publishTimes)), zap.Error(err))
		return
	}
	// 写入订单
	rec := recharge.UserRechargeRecord{
		BotID:           botID,
		PublishTimes:    publishTimes,
		StartTime:       time.Now(),
		PublishInterval: 30,
		PublishContent:  val,
		Status:          global.AdRechargeCreate, // 创建
		UserID:          userID,
		UpdateID:        updateID,
		Price:           price,
	}
	if err = global.GVA_DB.Create(&rec).Error; err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("rec", rec), zap.Error(err))
		return
	}

	if err = global.GVA_REDIS.Del(ctx, draftKey).Err(); err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.String("draftKey", draftKey), zap.Error(err))
		return
	}

	sendTex := fmt.Sprintf("✅ 广告订单创建成功，请前往后台完成支付, 支付价格为: %f", price)
	botHandler.SendTextMessage(chatID, token, sendTex)

	// 查询所有频道
	var channels []bot.BotChannel
	if err = global.GVA_DB.Where("bot_id = ?", botID).Find(&channels).Error; err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Error(err))
		return
	}

	var medias []bot_handler.MediaItem
	if err = json.Unmarshal([]byte(val), &medias); err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("val", val), zap.Error(err))
		return
	}

	cmdCfg := cache.NewBotCmdCache(botID, global.BotReplyCnfPublish2Channel, global.BotReplyCnfType)
	if _, err = cache.CacheGetItem(cmdCfg); err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Error(err))
		return
	}

	buttons := ParseContentFromCfg(*cmdCfg, global.ButtonTypeInline)
	global.GVA_LOG.Debug("botHandle HandleAdConfirm", zap.Any("buttons", buttons))
	// 不管了，都发吧
	for _, ch := range channels {
		if _, err = botHandler.TgSend(token, ch.ChannelID, medias, buttons); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm TgSend", zap.Int64("channelID", ch.ChannelID), zap.Error(err))
		}
	}
	return nil
}
