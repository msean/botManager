package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
)

func ReceiveAdContentHandle(update tgbotapi.Update, token string, botID int64) {
	msg := update.Message
	if msg == nil {
		return
	}

	medias := ParseIncomingMedia(msg)
	if len(medias) == 0 {
		return
	}

	ctx := context.Background()
	updateID := update.UpdateID // 每条内容唯一 ID
	userID := getChatUserID(update)

	data, _ := json.Marshal(medias)

	// 设置 30 分钟过期（自动处理超时）
	global.GVA_REDIS.Set(ctx, cache.AdDraftCacheKey(botID, userID, int64(update.UpdateID)), string(data), confirmAdExpire)

	// 清除等待状态
	global.GVA_REDIS.Del(ctx, cache.AdWaitCacheKey(botID, userID))

	buttons := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("✅ 确认发布", fmt.Sprintf("%s:%d", AdConfirmCmd, updateID)),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("❌ 取消发布", fmt.Sprintf("%s:%d", AdCancelCmd, updateID)),
		},
	)
	// 一次性发送预览 + 按 updateID 绑定按钮
	bot_handler.TgSend(token, msg.Chat.ID, medias, buttons)
}

func ParseIncomingMedia(msg *tgbotapi.Message) []bot_handler.MediaItem {
	var items []bot_handler.MediaItem

	if msg.Caption != "" {
		items = append(items, bot_handler.MediaItem{Type: "text", Text: msg.Caption})
	}
	if msg.Text != "" {
		items = append(items, bot_handler.MediaItem{Type: "text", Text: msg.Text})
	}

	if len(msg.Photo) > 0 {
		ph := msg.Photo[len(msg.Photo)-1]
		items = append(items, bot_handler.MediaItem{Type: "photo", FileID: ph.FileID})
	}

	if msg.Video != nil {
		items = append(items, bot_handler.MediaItem{Type: "video", FileID: msg.Video.FileID})
	}

	return items
}

func SendPreviewWithButtons(chatID int64, token string, medias []bot_handler.MediaItem, updateID int) {
	// 创建按钮

}
