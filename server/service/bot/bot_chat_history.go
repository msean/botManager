package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/service/cache"
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

		if Aerr := global.GVA_PGSQL.Exec(fmt.Sprintf(
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_%d_chat_id ON "%s"(id)`,
			svc.chatGroupID, svc.MessageTable(),
		)).Error; Aerr != nil {
			global.GVA_LOG.Error("create index chat_id err", zap.Error(Aerr))
		}
		if Berr := global.GVA_PGSQL.Exec(
			fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%d_time ON %s (timestamp)`,
				svc.chatGroupID, svc.MessageTable()),
		).Error; Berr != nil {
			global.GVA_LOG.Error("create index time err", zap.Error(Berr))
		}
		if Cerr := global.GVA_PGSQL.Exec(
			fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%d_user_id ON %s (user_id)`,
				svc.chatGroupID, svc.MessageTable()),
		).Error; Cerr != nil {
			global.GVA_LOG.Error("create index user", zap.Error(Cerr))
		}
		if Derr := global.GVA_PGSQL.Exec(fmt.Sprintf(
			`CREATE INDEX IF NOT EXISTS idx_%d_fts ON "%s"
			 USING gin (to_tsvector('simple', coalesce(text,'') || ' ' || coalesce(caption,'')))`,
			svc.chatGroupID, svc.MessageTable(),
		)).Error; Derr != nil {
			global.GVA_LOG.Error("create index user", zap.Error(Derr))
		}
		if Ferr := global.GVA_PGSQL.Exec(fmt.Sprintf(
			`CREATE INDEX IF NOT EXISTS idx_%d_user_name ON %s (lower(username))`,
			svc.chatGroupID, svc.MessageTable(),
		)).Error; Ferr != nil {
			global.GVA_LOG.Error("create index user", zap.Error(Ferr))
		}
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

func GetTelegramFileURL(botToken, fileID string) (string, error) {
	apiURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getFile?file_id=%s",
		botToken,
		fileID,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if !result.OK || result.Result.FilePath == "" {
		return "", errors.New("getFile failed")
	}

	fileURL := fmt.Sprintf(
		"https://api.telegram.org/file/bot%s/%s",
		botToken,
		result.Result.FilePath,
	)

	return fileURL, nil
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
		record.NickName = fmt.Sprintf("%s%s", msg.From.FirstName, msg.From.LastName)
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

func (svc BotChatHistorySvc) QueryMessages(
	req request.ChatMessageQuery,
) (
	list []bot.TgChatMessageV1,
	hasMore bool,
	err error,
) {

	botCache := cache.NewBotCache(req.BotID)
	var has bool
	if has, err = cache.CacheGetItem(botCache); !has || err != nil {
		return
	}
	table := `"` + svc.MessageTable() + `" m`

	db := global.GVA_PGSQL.
		Table(table).
		Select(`
			m.*,
			r.id as reply_id,
			r.user_id as reply_user_id,
			r.username as reply_username,
			r.text as reply_text,
			r.message_type as reply_message_type
		`).
		Joins(`LEFT JOIN "` + svc.MessageTable() + `" r ON m.reply_to_message_id = r.id`)

	// ===== 查询条件 =====
	if req.UserID > 0 {
		db = db.Where("m.user_id = ?", req.UserID)
	}
	if req.Username != "" {
		db = db.Where("lower(m.username) LIKE ?", "%"+strings.ToLower(req.Username)+"%")
	}
	// if req.Text != "" {

	// 	db = db.Where("(m.text ILIKE ? OR m.caption ILIKE ?)", "%"+req.Text+"%", "%"+req.Text+"%")
	// }
	if req.Text != "" {
		db = db.Where(`
			to_tsvector(
				'simple',
				coalesce(m.text,'') || ' ' || coalesce(m.caption,'')
			) @@ plainto_tsquery('simple', ?)
		`, req.Text)
	}
	if req.StartTime != nil {
		db = db.Where("m.timestamp >= ?", *req.StartTime)
	}
	if req.EndTime != nil {
		db = db.Where("m.timestamp <= ?", *req.EndTime)
	}

	limit := req.Limit + 1 // 多查一条判断是否还有

	// ===== 游标分页 =====
	switch {
	case req.BeforeID > 0:
		// 加载更早
		db = db.
			Where("m.id < ?", req.BeforeID).
			Order("m.id DESC")

	case req.AfterID > 0:
		// 加载更新
		db = db.
			Where("m.id > ?", req.AfterID).
			Order("m.id ASC")

	default:
		// 初始化：最早 → 最新
		db = db.Order("m.id ASC")
	}

	var rows []bot.TgChatMessageV1
	if err = db.Limit(limit).Find(&rows).Error; err != nil {
		return
	}

	for i := range rows {
		msg := &rows[i]
		if msg.FileID != "" && (msg.FileType == "photo" || msg.FileType == "video") {
			url, err := GetTelegramFileURL(botCache.Token, msg.FileID)
			if err == nil {
				msg.FileUrl = &url
			}
		}
	}

	if len(rows) > req.Limit {
		hasMore = true
		rows = rows[:req.Limit]
	}

	// ⚠️ 只有 beforeID 需要反转
	if req.BeforeID > 0 {
		for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
			rows[i], rows[j] = rows[j], rows[i]
		}
	}

	list = rows
	return
}
