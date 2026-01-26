package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/bot"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func HandleAdCancel(update botapi.Update, token string, botID int64) (err error) {
	userID := bot_handler.GetChatUserID(update)
	chatID := bot_handler.GetChatID(update)
	ctx := context.Background()

	data := update.CallbackQuery.Data
	parts := strings.Split(data, ":")
	if len(parts) == 1 {
		return
	}
	draftMsgIDStr := parts[1]
	draftMsgID, _ := strconv.Atoi(draftMsgIDStr)
	draftKey := cache.AdDraftCacheKey(botID, userID, draftMsgID)
	bot, _ := botapi.NewBotAPI(token)

	// // 2. 判断草稿是否存在
	// val, _ := global.GVA_REDIS.Get(ctx, draftKey).Result()

	// if val == "" {
	// 	// 草稿不存在 = 超时
	// 	bot.Send(botapi.NewMessage(chatID,
	// 		"⏱️ 发布请求已超时，请重新提交内容。"))
	// 	return nil
	// }
	// if err = dao.RechargeDao.CancelOrder(global.GVA_MYSQL, botID, userID, msgID); err != nil {
	// 	global.GVA_LOG.Error("HandleAdCancel CancelOrder", zap.Error(err))
	// 	return
	// }

	// 3. 正常取消
	global.GVA_REDIS.Del(ctx, draftKey)

	del := botapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID)
	bot.Send(del)

	bot.Send(botapi.NewMessage(chatID, "❌ 已取消发布。"))

	return nil
}

// 确认发布
func HandleAdConfirmCallback(update botapi.Update, token string, botID int64) (err error) {
	publishTimes := 1
	userID := bot_handler.GetChatUserID(update)
	chatID := bot_handler.GetChatID(update)
	userName := bot_handler.GetUserName(update)
	ctx := context.Background()

	data := update.CallbackQuery.Data
	// userName := update.CallbackQuery.From.UserName
	// if userName == "" {
	// 	userName = update.CallbackQuery.From.FirstName + " " + update.CallbackQuery.From.LastName
	// }
	parts := strings.Split(data, ":")
	if len(parts) == 1 {
		return
	}
	draftMsgIDStr := parts[1]
	draftMsgID, _ := strconv.Atoi(draftMsgIDStr)
	draftKey := cache.AdDraftCacheKey(botID, userID, draftMsgID)
	var botHandler *bot_handler.Bot
	if botHandler, err = bot_handler.NewBot(token); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(draftMsgID)), zap.Error(err))
		return
	}

	val, err := global.GVA_REDIS.Get(ctx, draftKey).Result()
	if err != nil || val == "" {
		if err = botHandler.DeleteMsg(chatID, draftMsgID); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm DeleteMsg", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(draftMsgID)), zap.Error(err))
		}
		if err = botHandler.SendTextMessage(chatID, "❌ 此发布请求已过期，请重新发送内容。"); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.String("userName", userName), zap.Int64("msgID", int64(draftMsgID)), zap.Error(err))
		}
		return nil
	}

	var wallet recharge.UserWallet
	if wallet, err = dao.RechargeDao.GetUserWallet(global.GVA_MYSQL, botID, userID, userName); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm GetUserWallet", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.String("userName", userName), zap.Int64("userID", userID), zap.Int64("msgID", int64(draftMsgID)), zap.Error(err))
		if err = botHandler.SendTextMessage(chatID, "获取余额失败，稍后再试"); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(draftMsgID)), zap.Error(err))
			return
		}
	}

	var cnf cache.RechargeCnfObj
	var has bool
	if cnf, has, err = cache.NewRechargeCnfListCache(botID).WherePublishTimes(publishTimes); !has || err != nil {
		global.GVA_LOG.Error("HandleAdConfirm RechargeCnfListCache", zap.Int64("botID", botID),
			zap.Int64("chatID", chatID),
			zap.String("userName", userName),
			zap.Int64("userID", userID),
			zap.Int64("msgID", int64(draftMsgID)),
			zap.Bool("has", has),
			zap.Error(err),
		)
		if err = botHandler.SendTextMessage(chatID, "后台价格配置有误，稍后再试"); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(draftMsgID)), zap.Error(err))
			return
		}
	}
	global.GVA_LOG.Debug("HandleAdConfirm recharge", zap.Int64("botID", botID), zap.Any("cnf", cnf), zap.Any("wallet.Balance", wallet.Balance), zap.Any("publishTimes", publishTimes))

	// 余额不足提示充值
	if wallet.Balance < cnf.Price {
		msg := botapi.NewMessage(chatID, fmt.Sprintf("你当前的余额为%.3f, 余额不足，请充值", wallet.Balance))
		btn := botapi.NewInlineKeyboardButtonData("⚡ 立即充值", NoticeRechargeCmd)
		keyboard := botapi.NewInlineKeyboardMarkup(
			botapi.NewInlineKeyboardRow(btn),
		)

		msg.ReplyMarkup = keyboard
		if _, err = botHandler.Send(msg); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm SendTextMessage", zap.Int64("botID", botID), zap.Int64("chatID", chatID), zap.Int64("msgID", int64(draftMsgID)), zap.Error(err))
			return
		}
		global.GVA_REDIS.Set(ctx, cache.AdDraftConfirmCacheKey(botID, userID), val, constant.OrderMatchAgo*time.Minute)
	}

	var medias []bot_handler.MediaItem
	if err = json.Unmarshal([]byte(val), &medias); err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("val", val), zap.Error(err))
		return
	}

	// 余额充足 立马 扣减余额
	if _, err = dao.RechargeDao.ReduceBalance(global.GVA_MYSQL, botID, userID, cnf.Price); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm ReduceBalance", zap.Int64("botID", botID), zap.Int64("userID", userID), zap.Any("price", cnf.Price), zap.Error(err))
		return
	}

	hook := func(channels []cache.BotChannelCache) error {
		go func() {
			var channelIDList []int64
			for _, channel := range channels {
				channelIDList = append(channelIDList, channel.ChannelID)
			}
			err := dao.RechargeDao.CreatePublishRecords(
				global.GVA_MYSQL,
				recharge.AdPublishRecord{
					BotID:        botID,
					PublishTimes: 1,
					UserID:       userID,
					UserName:     userName,
					Price:        cnf.Price,
					Content:      val,
				},
				channelIDList,
			)
			if err != nil {
				global.GVA_LOG.Error("保存发布记录失败", zap.Error(err))
			}
		}()
		return nil
	}
	// 发布到所有渠道
	if err = bot.NewBotHandlerSvc(botID).PublishAd2Channel(*botHandler, chatID, medias, hook); err != nil {
		global.GVA_LOG.Error("botHandle PublishAd2Channel", zap.Int("botID", int(botID)), zap.Any("val", val), zap.Error(err))
		return
	}

	return
}
