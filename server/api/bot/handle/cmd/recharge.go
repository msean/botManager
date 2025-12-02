package cmd

import (
	"context"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/service/recharge"
	"github.com/msean/botmanager/server/utils/bot_handler"
)

func Recharge(chatID int64, userID int64, updateID int, token string, botID int64, msgID int, amount float64) {
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
		pay.Recharge(token, userID, chatID, int64(updateID), amount)
	}
}
