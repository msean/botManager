package private

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	rechargeModel "github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/service/recharge"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Recharge(chatID int64, userID int64, token string, botID int64, msgID int, userName string, amount float64) (err error) {
	// 自定义
	if amount == 0 {
		// 设置当前用户状态
		global.GVA_REDIS.Set(context.Background(), cache.RechargeTryCountKey(botID, userID), 0, waitAdContentExpire)
		global.GVA_REDIS.Set(context.Background(), cache.AdWaitCacheKey(botID, userID), waitRechargeState, waitAdContentExpire)
		botApi, _ := bot_handler.NewBot(token)
		botApi.SendTextMessage(chatID, "请输入充值金额，⚠️注意：输入必须为数字，默认单位是USDT")
	} else {
		// 创建订单
		pay := recharge.NewPay(botID)
		if err = pay.Recharge(token, userID, chatID, msgID, userName, amount); err != nil {
			global.GVA_LOG.Error("Recharge", zap.Int64("botID", botID), zap.Int64("userID", userID), zap.Error(err))
		}
	}
	return
}

func RechargeChoiceHandler(update botapi.Update, token string, botID int64) (err error) {
	ctx := context.Background()
	key := cache.RechargeTryCountKey(botID, update.Message.From.ID)
	// 当前次数
	tryCount, _ := global.GVA_REDIS.Get(ctx, key).Int()

	if tryCount >= 1 {
		// 清理状态
		global.GVA_REDIS.Del(ctx,
			cache.AdWaitCacheKey(botID, update.Message.From.ID),
			key,
		)

		reply := botapi.NewMessage(
			update.Message.Chat.ID,
			"❌ 输入错误次数过多，请重新点击充值按钮",
		)
		bot, _ := botapi.NewBotAPI(token)
		bot.Send(reply)
		return
	}

	text := strings.TrimSpace(update.Message.Text)

	amount, err := strconv.ParseFloat(text, 64)
	if err != nil || amount <= 0 {

		// 次数 +1
		global.GVA_REDIS.Incr(ctx, key)

		left := 2 - (tryCount + 1)

		reply := botapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("❌ 金额无效(输入必须为数字，默认单位是USDT)，你还有 %d 次输入机会", left),
		)

		bot, _ := botapi.NewBotAPI(token)
		bot.Send(reply)
		return
	}

	// ✅ 输入正确 -> 清理状态
	global.GVA_REDIS.Del(ctx,
		cache.AdWaitCacheKey(botID, update.Message.From.ID),
		key,
	)
	Recharge(update.Message.Chat.ID, update.Message.From.ID, token, botID, update.Message.MessageID, bot_handler.GetUserName(update), amount)
	return
}

func RechargeCancelHandler(update botapi.Update, token string, botID int64) (err error) {
	data := update.CallbackQuery.Data
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return
	}

	cancelRechargeID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		global.GVA_LOG.Error(
			"RechargeCancelHandler parse error",
			zap.String("data", data),
			zap.Error(err),
		)
		return
	}

	if cancelRechargeID <= 0 {
		return nil
	}

	chatID := update.CallbackQuery.Message.Chat.ID
	msgID := update.CallbackQuery.Message.MessageID

	bot, _ := botapi.NewBotAPI(token)

	// 修改数据库记录为【已取消】
	err = global.GVA_MYSQL.Model(&rechargeModel.UserRechargeRecord{}).
		Where("id = ?", cancelRechargeID).
		Update("status", constant.AdRechargeCancel).Error // 4 = 已取消(你自定义)
	if err != nil {
		bot.Send(botapi.NewMessage(chatID, "❌ 取消失败，请稍后重试"))
		return
	}

	// 删除原消息
	bot.Send(botapi.NewDeleteMessage(chatID, msgID))

	// 提示用户
	bot.Send(botapi.NewMessage(chatID, "✅ 充值已取消"))

	return nil
}

func RechargeInputCallbackHandler(update botapi.Update, token string, botID int64) (err error) {
	data := update.CallbackQuery.Data
	parts := strings.Split(data, "_")
	if len(parts) == 1 {
		return
	}

	_amount := parts[1]
	chatID := bot_handler.GetChatID(update)
	userID := bot_handler.GetChatUserID(update)
	msgID := update.CallbackQuery.Message.MessageID

	if data == "/rechargeChoice_close" {
		botApi, _ := bot_handler.NewBot(token)
		return botApi.DeleteOriginMessage(chatID, msgID)
	}
	global.GVA_LOG.Debug("BotMsgHandlerSvc RechargeCallbackHandler", zap.Any("msgID", msgID), zap.Any("data", data))
	var amount float64
	if amount, err = strconv.ParseFloat(_amount, 64); err != nil {
		global.GVA_LOG.Error(
			"amount parse failed",
			zap.String("data", data),
			zap.Error(err),
		)
		return
	}

	global.GVA_LOG.Debug("BotMsgHandlerSvc CallbackQuery amount", zap.Any("amount", amount), zap.Any("msgID", msgID), zap.Any("data", data))
	Recharge(chatID, userID, token, botID, msgID, bot_handler.GetUserName(update), amount)
	botApi, _ := bot_handler.NewBot(token)
	return botApi.DeleteOriginMessage(chatID, msgID)
}

func CheckRechargeAndNotify(botID, userID, chatID int64, token string) (canRecharge bool, err error) {
	var record rechargeModel.UserRechargeRecord

	err = global.GVA_MYSQL.
		Where("bot_id = ? AND user_id = ? AND status = ?", botID, userID, constant.AdRechargeCreate).
		Order("id DESC").
		First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}

		global.GVA_LOG.Error("CheckRechargeAndNotify error", zap.Error(err))
		return false, err
	}

	msg := fmt.Sprintf(
		"⚠️ 你当前有一笔进行中的充值订单\n\n💰 金额：%.3f USDT\n📍 地址：%s\n🕒 创建时间：%s\n\n请先完成或取消该订单，再发起新的充值。",
		record.Price,
		record.PaymentAddr,
		record.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	// // 取消按钮
	// btnCancel := botapi.NewInlineKeyboardButtonData(
	// 	"❌ 取消当前订单",
	// 	fmt.Sprintf("recharge_cancel:%d", record.ID),
	// )

	// // 显示金额按钮（可选）
	// btnAmount := botapi.NewInlineKeyboardButtonData(
	// 	fmt.Sprintf("💰 %.3f USDT", record.Price),
	// 	"amount_show",
	// )

	// keyboard := botapi.NewInlineKeyboardMarkup(
	// 	[]botapi.InlineKeyboardButton{btnAmount},
	// 	[]botapi.InlineKeyboardButton{btnCancel},
	// )

	bot, err := botapi.NewBotAPI(token)
	if err != nil {
		global.GVA_LOG.Error("CheckRechargeAndNotify NewBotAPI", zap.Error(err))
		return false, err
	}

	send := botapi.NewMessage(chatID, msg)
	// send.ReplyMarkup = keyboard

	if _, err = bot.Send(send); err != nil {
		global.GVA_LOG.Error("send message error", zap.Error(err))
		return false, err
	}

	return false, nil
}
