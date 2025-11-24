package cmd

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
)

func PublishAdHandle(update tgbotapi.Update, token string, botID int64) {
	var userID int64
	if update.CallbackQuery != nil {
		userID = int64(update.CallbackQuery.From.ID)
	} else if update.Message != nil {
		userID = int64(update.Message.From.ID)
	}
	// 设置用户状态：等待输入内容
	cacheKey := fmt.Sprintf("bot_manager:bot:%d:user:%d:state", botID, userID)
	global.GVA_REDIS.Set(context.Background(), cacheKey, waitAdContentState, waitAdContentExpire)
}
