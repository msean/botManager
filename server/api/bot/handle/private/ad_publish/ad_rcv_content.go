package adpublish

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func ReceiveAdContentHandle(update botapi.Update, token string, botID int64) (err error) {
	msg := update.Message
	if msg == nil {
		return
	}
	msgID := msg.MessageID

	medias := ParseIncomingMedia(msg)
	if len(medias) == 0 {
		return
	}

	var botHandler *bot_handler.Bot
	if botHandler, err = bot_handler.NewBot(token); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", botID), zap.Error(err))
		return
	}

	ctx := context.Background()
	userID := bot_handler.GetChatUserID(update)

	data, _ := json.Marshal(medias)

	// 设置 30 分钟过期（自动处理超时）
	global.GVA_REDIS.Set(ctx, cache.AdDraftCacheKey(botID, userID, msgID), string(data), confirmAdExpire)

	// 清除等待状态
	global.GVA_REDIS.Del(ctx, cache.AdWaitCacheKey(botID, userID))

	buttons := botapi.NewInlineKeyboardMarkup(
		[]botapi.InlineKeyboardButton{
			botapi.NewInlineKeyboardButtonData("✅ 确认发布", AdConfirmCmd+":"+strconv.Itoa(msgID)),
		},
		[]botapi.InlineKeyboardButton{
			botapi.NewInlineKeyboardButtonData("❌ 取消发布", AdCancelCmd+":"+strconv.Itoa(msgID)),
		},
	)
	// 一次性发送预览 + 按 updateID 绑定按钮
	if _, err = botHandler.TgSend(msg.Chat.ID, medias, buttons); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", botID), zap.Any("medias", medias), zap.Any("buttons", buttons), zap.Error(err))
	}
	return
}

func ParseIncomingMedia(msg *botapi.Message) []bot_handler.MediaItem {
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
