package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/service/cache"
)

type MediaItem struct {
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
	FileID string `json:"file_id,omitempty"`
}

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

	// 一次性发送预览 + 按 updateID 绑定按钮
	SendPreviewWithButtons(msg.Chat.ID, token, medias, updateID)
}

func ParseIncomingMedia(msg *tgbotapi.Message) []MediaItem {
	var items []MediaItem

	if msg.Caption != "" {
		items = append(items, MediaItem{Type: "text", Text: msg.Caption})
	}
	if msg.Text != "" {
		items = append(items, MediaItem{Type: "text", Text: msg.Text})
	}

	if len(msg.Photo) > 0 {
		ph := msg.Photo[len(msg.Photo)-1]
		items = append(items, MediaItem{Type: "photo", FileID: ph.FileID})
	}

	if msg.Video != nil {
		items = append(items, MediaItem{Type: "video", FileID: msg.Video.FileID})
	}

	return items
}
func SendPreviewWithButtons(chatID int64, token string, medias []MediaItem, updateID int) {
	bot, _ := tgbotapi.NewBotAPI(token)

	// 先按顺序发内容
	for _, m := range medias {
		switch m.Type {
		case "text":
			msg := tgbotapi.NewMessage(chatID, m.Text)
			bot.Send(msg)

		case "photo":
			msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(m.FileID))
			bot.Send(msg)

		case "video":
			msg := tgbotapi.NewVideo(chatID, tgbotapi.FileID(m.FileID))
			bot.Send(msg)
		}
	}

	// 最后一条带按钮（包含 updateID）
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("✅ 确认发布", fmt.Sprintf("AdConfirm:%d", updateID)),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("❌ 取消发布", fmt.Sprintf("AdCancel:%d", updateID)),
		},
	)

	msg := tgbotapi.NewMessage(chatID, "请确认是否发布此广告：")
	msg.ReplyMarkup = buttons
	bot.Send(msg)
}
