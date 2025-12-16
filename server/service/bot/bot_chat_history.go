package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/bot/request"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

type BotChatHistorySvc struct {
	botID       int64
	chatGroupID int64
}

func NewBotChatHistorySvc(botID, chatGroupID int64) *BotChatHistorySvc {
	return &BotChatHistorySvc{
		botID:       botID,
		chatGroupID: chatGroupID,
	}
}

func (svc BotChatHistorySvc) MessageTable() string {
	return fmt.Sprintf("bot_messages_%d_%d", svc.botID, svc.chatGroupID)
}

func (svc BotChatHistorySvc) tableLockKey() string {
	return fmt.Sprintf("lock:tg_msg_table:%d", svc.chatGroupID)
}

func (svc BotChatHistorySvc) Sync() error {
	global.GVA_LOG.Debug("BotChatHistorySvc Sync", zap.String("table", svc.MessageTable()))
	ctx := context.Background()

	cacheKey := svc.tableCacheKey()
	lockKey := svc.tableLockKey()

	// 1️⃣ 快速路径：缓存已存在
	if ok, _ := global.GVA_REDIS.Get(ctx, cacheKey).Bool(); ok {
		return nil
	}

	global.GVA_LOG.Debug("BotChatHistorySvc Sync", zap.String("cacheKey", cacheKey))
	// 2️⃣ 尝试加分布式锁（10 秒足够）
	locked, err := global.GVA_REDIS.SetNX(ctx, lockKey, 1, 10*time.Second).Result()
	if err != nil {
		return err
	}
	if !locked {
		return nil
	}

	defer global.GVA_REDIS.Del(ctx, lockKey)

	if ok, _ := global.GVA_REDIS.Get(ctx, cacheKey).Bool(); ok {
		return nil
	}

	exists, err := svc.tableExists()
	if err != nil {
		return err
	}

	global.GVA_LOG.Debug("BotChatHistorySvc Sync", zap.Bool("exists", exists))
	if !exists {
		if err := global.GVA_PGSQL.
			Table(svc.MessageTable()).
			AutoMigrate(&bot.TGMessageRecord{}); err != nil {
			return err
		}

		global.GVA_PGSQL.Exec(
			fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%d_time ON %s (timestamp)`,
				svc.chatGroupID, svc.MessageTable()),
		)
		global.GVA_PGSQL.Exec(
			fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%d_user ON %s (user_id)`,
				svc.chatGroupID, svc.MessageTable()),
		)
	}

	// 7️⃣ 写缓存
	_ = global.GVA_REDIS.Set(
		ctx,
		cacheKey,
		1,
		24*time.Hour,
	).Err()

	return nil
}

func (svc BotChatHistorySvc) tableExists() (bool, error) {
	var exists bool
	sql := `
SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = 'public'
      AND table_name = ?
)
`
	err := global.GVA_PGSQL.Raw(sql, svc.MessageTable()).Scan(&exists).Error
	return exists, err
}

func (svc BotChatHistorySvc) tableCacheKey() string {
	return fmt.Sprintf("tg:msg_table:exists:%d:%d", svc.botID, svc.chatGroupID)
}

func (svc BotChatHistorySvc) SaveMessage(update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	msg := update.Message

	if msg.Chat == nil || (msg.Chat.Type != "group" && msg.Chat.Type != "supergroup") {
		return nil
	}

	chatID := msg.Chat.ID
	record := bot.TGMessageRecord{
		ID:        int64(msg.MessageID),
		ChatID:    chatID,
		UserID:    0,
		IsBot:     false,
		Text:      msg.Text,
		Caption:   msg.Caption,
		Timestamp: time.Unix(int64(msg.Date), 0),
		CreatedAt: time.Now(),
	}

	// ===== 用户信息 =====
	if msg.From != nil {
		record.UserID = msg.From.ID
		record.Username = msg.From.UserName
		record.FirstName = msg.From.FirstName
		record.LastName = msg.From.LastName
		record.IsBot = msg.From.IsBot
	}

	if msg.ReplyToMessage != nil {
		record.ReplyToMessageID = int64(msg.ReplyToMessage.MessageID)
	}

	switch {
	case msg.Text != "":
		record.MessageType = "text"

	case len(msg.Photo) > 0:
		record.MessageType = "photo"
		photo := msg.Photo[len(msg.Photo)-1]
		record.FileID = photo.FileID
		record.FileUniqueID = photo.FileUniqueID
		record.FileType = "photo"

	case msg.Video != nil:
		record.MessageType = "video"
		record.FileID = msg.Video.FileID
		record.FileUniqueID = msg.Video.FileUniqueID
		record.FileType = "video"

	case msg.Document != nil:
		record.MessageType = "document"
		record.FileID = msg.Document.FileID
		record.FileUniqueID = msg.Document.FileUniqueID
		record.FileType = msg.Document.MimeType

	case msg.Sticker != nil:
		record.MessageType = "sticker"
		record.FileID = msg.Sticker.FileID
		record.FileUniqueID = msg.Sticker.FileUniqueID
		record.FileType = "sticker"

	case msg.Voice != nil:
		record.MessageType = "voice"
		record.FileID = msg.Voice.FileID
		record.FileUniqueID = msg.Voice.FileUniqueID
		record.FileType = "voice"

	case msg.Audio != nil:
		record.MessageType = "audio"
		record.FileID = msg.Audio.FileID
		record.FileUniqueID = msg.Audio.FileUniqueID
		record.FileType = msg.Audio.MimeType

	default:
		record.MessageType = "unknown"
	}

	raw, _ := json.Marshal(msg)
	record.Raw = datatypes.JSON(raw)

	return global.GVA_PGSQL.
		Table(svc.MessageTable()).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).
		Create(&record).Error
}

func (svc BotChatHistorySvc) QueryMessages(req request.ChatMessageQuery) (list []bot.TgChatMessageV1, total int64, err error) {
	table := svc.MessageTable()

	quotedTable := `"` + table + `"`

	db := global.GVA_PGSQL.
		Table(quotedTable + " m").
		Select(`
			m.*,
			r.id as reply_id,
			r.user_id as reply_user_id,
			r.username as reply_username,
			r.text as reply_text,
			r.message_type as reply_message_type
		`).
		Joins(
			"LEFT JOIN " + quotedTable + " r ON m.reply_to_message_id = r.id",
		)

	if req.UserID != 0 {
		db = db.Where("m.user_id = ?", req.UserID)
	}
	if req.Username != "" {
		db = db.Where("m.username ILIKE ?", "%"+req.Username+"%")
	}
	if !req.StartTime.IsZero() {
		db = db.Where("m.timestamp >= ?", req.StartTime)
	}

	if !req.EndTime.IsZero() {
		db = db.Where("m.timestamp <= ?", req.EndTime)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	err = db.
		Order("timestamp asc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&list).Error

	return
}
