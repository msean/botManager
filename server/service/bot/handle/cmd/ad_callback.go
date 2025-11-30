package cmd

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	rechargeSrv "github.com/msean/botmanager/server/service/recharge"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

func HandleAdCancel(chatID int64, userID int64, updateID int, token string, botID int64, msgID int) (err error) {
	ctx := context.Background()
	bot, _ := tgbotapi.NewBotAPI(token)

	// // 2. 判断草稿是否存在
	draftKey := cache.AdDraftCacheKey(botID, userID, int64(updateID))
	// val, _ := global.GVA_REDIS.Get(ctx, draftKey).Result()

	// if val == "" {
	// 	// 草稿不存在 = 超时
	// 	bot.Send(tgbotapi.NewMessage(chatID,
	// 		"⏱️ 发布请求已超时，请重新提交内容。"))
	// 	return nil
	// }

	if err = dao.RechargeDao.CancelOrder(global.GVA_DB, botID, userID, int64(updateID)); err != nil {
		global.GVA_LOG.Error("HandleAdCancel CancelOrder", zap.Error(err))
		return
	}

	// 3. 正常取消
	global.GVA_REDIS.Del(ctx, draftKey)

	del := tgbotapi.NewDeleteMessage(chatID, msgID)
	bot.Send(del)

	bot.Send(tgbotapi.NewMessage(chatID, "❌ 已取消发布。"))

	return nil
}

func HandleAdConfirm(chatID int64, userID int64, userName string, updateID int64, token string, botID int64, msgID int, publishTimes int) (err error) {
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

	// 获取价格
	var price float64
	pay := rechargeSrv.NewPay(botID, publishTimes)
	if price, err = pay.RandomPrice(); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm RandomPrice", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("publishTimes", int64(publishTimes)), zap.Error(err))
		return
	}
	// 获取收款地址
	var paymentAddr string
	if paymentAddr, err = pay.GetPaymentAddr(); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm paymentAddr", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Error(err))
		return
	}

	// 写入订单
	rec := recharge.UserRechargeRecord{
		BotID:           botID,
		PublishTimes:    publishTimes,
		StartTime:       time.Now(),
		PublishInterval: 30,
		PublishContent:  val,
		Status:          constant.AdRechargeCreate, // 创建
		UserID:          userID,
		UpdateID:        updateID,
		Price:           price,
		PaymentAddr:     paymentAddr,
		UserName:        userName,
		ChatID:          chatID,
	}
	if err = global.GVA_DB.Create(&rec).Error; err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("rec", rec), zap.Error(err))
		return
	}

	// amount := fmt.Sprintf("%.3f", price)
	sendText := bot_handler.FormatRechargeMessage(
		rec.ID,
		fmt.Sprintf("%.3f", price),
		paymentAddr,
		rec.CreatedAt.Format("2006/01/02 15:04:05"),
		constant.OrderLeftPaid,
	)

	if err = botHandler.SendMarkDownMessage(chatID, token, sendText); err != nil {
		global.GVA_LOG.Error("botHandle SendMarkDownMessage", zap.Int("botID", int(botID)), zap.String("draftKey", draftKey), zap.Any("sendText", sendText), zap.Error(err))
		return
	}

	if err = global.GVA_REDIS.Del(ctx, draftKey).Err(); err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.String("draftKey", draftKey), zap.Error(err))
		return
	}
	// 查询所有频道
	// var medias []bot_handler.MediaItem
	// if err = json.Unmarshal([]byte(val), &medias); err != nil {
	// 	global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("val", val), zap.Error(err))
	// 	return
	// }

	// cmdCfg := cache.NewBotCmdCache(botID, global.BotReplyCnfPublish2Channel, global.BotReplyCnfType)
	// if _, err = cache.CacheGetItem(cmdCfg); err != nil {
	// 	global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Error(err))
	// 	return
	// }

	// buttons := ParseContentFromCfg(*cmdCfg, global.ButtonTypeInline)
	// global.GVA_LOG.Debug("botHandle HandleAdConfirm", zap.Any("buttons", buttons))
	return
}
