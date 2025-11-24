package cmd

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
)

func CancelPublishAdHandle(update tgbotapi.Update, token string, botID int64) {
	var userID int64
	if update.CallbackQuery != nil {
		userID = int64(update.CallbackQuery.From.ID)
	} else if update.Message != nil {
		userID = int64(update.Message.From.ID)
	}
	cacheKey := fmt.Sprintf("bot:%d:user:%d:state", botID, userID)
	global.GVA_REDIS.Del(context.Background(), cacheKey)
}
