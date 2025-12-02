package cmd

import (
	"context"
	"encoding/json"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/bot"
	"github.com/msean/botmanager/server/service/cache"
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

// 确认发布
func HandleAdConfirm(chatID int64, userID int64, userName string, updateID int64, token string, botID int64, msgID int, publishTimes int) (err error) {

	ctx := context.Background()

	draftKey := cache.AdDraftCacheKey(botID, userID, updateID)
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
		if err = botHandler.SendTextMessage(chatID, "❌ 此发布请求已过期，请重新发送内容。"); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.String("userName", userName), zap.Int64("msgID", int64(msgID)), zap.Error(err))
		}
		return nil
	}

	var wallet recharge.UserWallet
	if wallet, err = dao.RechargeDao.GetUserWallet(global.GVA_DB, botID, userID, userName); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm GetUserWallet", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.String("userName", userName), zap.Int64("userID", userID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
		if err = botHandler.SendTextMessage(chatID, "获取余额失败，稍后再试"); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
			return
		}
	}

	rechargeCnfList := cache.NewRechargeCnfListCache(botID)
	if _, err = cache.CacheGetItem(rechargeCnfList); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm RechargeCnfListCache", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.String("userName", userName), zap.Int64("userID", userID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
		if err = botHandler.SendTextMessage(chatID, "获取价格配置错误，稍后再试"); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
			return
		}
		return
	}
	cnf, has := rechargeCnfList.WherePublishTimes(publishTimes)
	if !has {
		global.GVA_LOG.Error("HandleAdConfirm RechargeCnfListCache", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.String("userName", userName), zap.Int64("userID", userID), zap.Int64("msgID", int64(msgID)))
		if err = botHandler.SendTextMessage(chatID, "后台价格配置有误，稍后再试"); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
			return
		}
	}

	global.GVA_LOG.Debug("HandleAdConfirm recharge", zap.Int64("botID", botID), zap.Any("cnf", cnf), zap.Any("wallet.Balance", wallet.Balance), zap.Any("publishTimes", publishTimes))

	// 余额不足提示充值
	if wallet.Balance < cnf.Price {
		msg := tgbotapi.NewMessage(chatID, "余额不足，请充值")
		btn := tgbotapi.NewInlineKeyboardButtonData("⚡ 立即充值", NoticeRechargeCmd)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(btn),
		)

		msg.ReplyMarkup = keyboard
		if _, err = botHandler.Send(msg); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(msgID)), zap.Error(err))
			return
		}
		global.GVA_REDIS.Set(ctx, cache.AdDraftConfirmCacheKey(botID, userID, updateID), val, constant.OrderMatchAgo*time.Minute)
	}

	var medias []bot_handler.MediaItem
	if err = json.Unmarshal([]byte(val), &medias); err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("val", val), zap.Error(err))
		return
	}

	// 余额充足 立马 扣减余额
	if _, err = dao.RechargeDao.ReduceBalance(global.GVA_DB, botID, userID, cnf.Price); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm ReduceBalance", zap.Int64("botID", botID), zap.Int64("userID", userID), zap.Any("price", cnf.Price), zap.Error(err))
		return
	}
	// 发布到所有渠道
	if err = bot.NewBotHandlerSvc(botID).PublishAd2Channel(*botHandler, chatID, medias); err != nil {
		global.GVA_LOG.Error("botHandle PublishAd2Channel", zap.Int("botID", int(botID)), zap.Any("val", val), zap.Error(err))
		return
	}
	return
}
