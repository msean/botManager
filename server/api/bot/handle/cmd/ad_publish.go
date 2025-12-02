package cmd

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/service/cache"
)

func PublishAdHandle(update tgbotapi.Update, token string, botID int64) {
	var userID int64
	if update.CallbackQuery != nil {
		userID = int64(update.CallbackQuery.From.ID)

	} else if update.Message != nil {
		userID = int64(update.Message.From.ID)

	}
	// 设置用户状态：等待输入内容
	global.GVA_REDIS.Set(context.Background(), cache.AdWaitCacheKey(botID, userID), waitAdContentState, waitAdContentExpire)
}

// func PublishAdCheckHandle(update tgbotapi.Update, token string, botID int64) (canPublic bool) {
// 	var userID int64
// 	if update.CallbackQuery != nil {
// 		userID = int64(update.CallbackQuery.From.ID)

// 	} else if update.Message != nil {
// 		userID = int64(update.Message.From.ID)

// 	}

// 	// 检查用户是否最近重复下单
// 	var has bool
// 	var err error
// 	if has, err = dao.RechargeDao.UserHasRecentOrder(global.GVA_DB, botID, userID); err != nil {
// 		global.GVA_LOG.Error("PublishAdHandle UserHasRecentOrder", zap.Int64("botID", botID), zap.Int64("userID", userID), zap.Error(err))
// 		return
// 	}
// 	if has {
// 		var botApi *bot_handler.Bot
// 		if botApi, err = bot_handler.NewBot(token); err != nil {
// 			global.GVA_LOG.Error("PublishAdHandle NewBot", zap.Int64("botID", botID), zap.Int64("userID", userID), zap.Error(err))
// 			return
// 		}
// 		botApi.SendTextMessage(update.Message.Chat.ID, "当前有未支付订单，若想重新下单，请先取消订单")
// 		return
// 	}
// 	return !has
// }
