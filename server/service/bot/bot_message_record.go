package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

type BotMsgRecordSvc struct {
	botID       int64
	chatGroupID int64
}

func NewBotMsgRecordSvc(botID, chatGroupID int64) *BotMsgRecordSvc {
	return &BotMsgRecordSvc{
		botID:       botID,
		chatGroupID: chatGroupID,
	}
}

func (svc BotMsgRecordSvc) MessageTable() string {
	return fmt.Sprintf("bot_messages_%d_%d", svc.botID, svc.chatGroupID)
}

func (svc BotMsgRecordSvc) tableLockKey() string {
	return fmt.Sprintf("lock:tg_msg_table:%d", svc.chatGroupID)
}

func (svc BotMsgRecordSvc) Sync() error {
	ctx := context.Background()

	cacheKey := svc.tableCacheKey()
	lockKey := svc.tableLockKey()

	// 1️⃣ 快速路径：缓存已存在
	if ok, _ := global.GVA_REDIS.Get(ctx, cacheKey).Bool(); ok {
		return nil
	}

	// 2️⃣ 尝试加分布式锁（10 秒足够）
	locked, err := global.GVA_REDIS.SetNX(ctx, lockKey, 1, 10*time.Second).Result()
	if err != nil {
		return err
	}
	if !locked {
		return nil
	}

	// 确保释放锁
	defer global.GVA_REDIS.Del(ctx, lockKey)

	// 3️⃣ 二次检查（非常关键）
	if ok, _ := global.GVA_REDIS.Get(ctx, cacheKey).Bool(); ok {
		return nil
	}

	// 4️⃣ 查表是否存在
	exists, err := svc.tableExists()
	if err != nil {
		return err
	}

	if !exists {
		// 5️⃣ 创建表
		if err := global.GVA_PGSQL.
			Table(svc.MessageTable()).
			AutoMigrate(&bot.TGMessageRecord{}); err != nil {
			return err
		}

		// 6️⃣ 创建索引
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

func (svc BotMsgRecordSvc) tableExists() (bool, error) {
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

func (svc BotMsgRecordSvc) tableCacheKey() string {
	return fmt.Sprintf("tg:msg_table:exists:%d:%d", svc.botID, svc.chatGroupID)
}

func (svc BotMsgRecordSvc) SaveMessage(update tgbotapi.Update) error {
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
