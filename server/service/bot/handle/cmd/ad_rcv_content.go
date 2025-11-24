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
	// 创建按钮
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("✅ 确认发布", fmt.Sprintf("AdConfirm:%d", updateID)),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("❌ 取消发布", fmt.Sprintf("AdCancel:%d", updateID)),
		},
	)

	// 调用通用发送函数
	AdSend(token, chatID, medias, buttons)
}

// SendAdMessage 统一发送广告内容（文字 / 图片 + 文本 / 视频 + 文本）
// 如果 replyMarkup != nil，则用作按钮，否则不带按钮
func AdSend(token string, chatID int64, medias []MediaItem, replyMarkup interface{}) (tgbotapi.Message, error) {
	var caption string
	var photoID string
	var videoID string

	bot, _ := tgbotapi.NewBotAPI(token)

	// 整理用户发送的数据
	for _, m := range medias {
		switch m.Type {
		case "text":
			caption = m.Text
		case "photo":
			photoID = m.FileID
		case "video":
			videoID = m.FileID
		}
	}

	// PHOTO
	if photoID != "" {
		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(photoID))
		msg.Caption = caption

		if replyMarkup != nil {
			msg.ReplyMarkup = replyMarkup
		}

		return bot.Send(msg)
	}

	// VIDEO
	if videoID != "" {
		msg := tgbotapi.NewVideo(chatID, tgbotapi.FileID(videoID))
		msg.Caption = caption

		if replyMarkup != nil {
			msg.ReplyMarkup = replyMarkup
		}

		return bot.Send(msg)
	}

	// ONLY TEXT
	msg := tgbotapi.NewMessage(chatID, caption)

	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}

	return bot.Send(msg)
}
