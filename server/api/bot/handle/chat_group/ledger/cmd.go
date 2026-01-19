package ledger

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	"github.com/msean/botmanager/server/service/cache"
	"go.uber.org/zap"
)

type (
	MessageParser interface {
		Match(botModel bot.Bot, update tgbotapi.Update) bool
		Handle() error
	}
	PermissionAware interface {
		NeedPermission() bool
		NeedAdmin() bool
		SetPermission(*cache.LedgerPermissionCache)
	}
)

type (
	ParserChain struct {
		parsers []MessageParser
	}
	ShouldPermissionAwareWithOutAdmin struct {
		confModel *cache.LedgerPermissionCache
	}
	NotPermissionAwareWithAdmin struct {
		confModel *cache.LedgerPermissionCache
	}
)

func (l *ShouldPermissionAwareWithOutAdmin) SetPermission(p *cache.LedgerPermissionCache) {
	l.confModel = p
}

func (l *ShouldPermissionAwareWithOutAdmin) NeedPermission() bool {
	return true
}

func (l *ShouldPermissionAwareWithOutAdmin) NeedAdmin() bool {
	return false
}

func (l *NotPermissionAwareWithAdmin) NeedPermission() bool {
	return false
}

func (l *NotPermissionAwareWithAdmin) NeedAdmin() bool {
	return true
}

func (l *NotPermissionAwareWithAdmin) SetPermission(p *cache.LedgerPermissionCache) {
	l.confModel = p
}
func (c *ParserChain) Handle(botModel bot.Bot, update tgbotapi.Update) error {

	for _, parser := range c.parsers {
		if !parser.Match(botModel, update) {
			continue
		}

		// 权限感知
		if p, ok := parser.(PermissionAware); ok {

			var isAdminUser bool
			isAdminUser, err := IsChatAdmin(
				botModel.Token,
				update.Message.Chat.ID,
				update.Message.From.ID,
			)
			if err != nil || !ok {
				return nil // 静默失败，不给回复
			}

			// 管理员权限（最高优先级）
			if p.NeedAdmin() || !isAdminUser {
				return nil
			}

			// 业务权限（Ledger 操作人）
			if p.NeedPermission() {
				permission, permit, err := HasPerMission(botModel, update, isAdminUser)
				if err != nil {
					return nil
				}
				if !isAdminUser && !permit {
					return nil
				}

				p.SetPermission(permission)
			}
		}

		return parser.Handle()
	}

	return nil
}

func NewParserChain(parsers ...MessageParser) *ParserChain {
	return &ParserChain{
		parsers: parsers,
	}
}

func Handle(botModel bot.Bot, tgMsg tgbotapi.Update) (err error) {
	if tgMsg.Message == nil || tgMsg.Message.Text == "" {
		return
	}
	chain := NewParserChain(
		&LedgerHandler{},
		&LedgerShowFeeRateHandler{},
		&LedgerSetFeeRateHandler{},
		&AddOprUserHandler{},
	)

	return chain.Handle(botModel, tgMsg)
}

func HasPerMission(botModel bot.Bot, tgMsg tgbotapi.Update, isAdmin bool) (ledgerPermission *cache.LedgerPermissionCache, permit bool, err error) {

	chatGroupID := tgMsg.Message.Chat.ID
	userID := tgMsg.Message.From.ID
	userName := tgMsg.Message.From.UserName

	ledgerPermission = cache.NewLedgerPermissionCache(botModel.BotID, chatGroupID)

	var has bool
	if has, err = cache.CacheGetItem(ledgerPermission); err != nil {
		global.GVA_LOG.Error("Ledger HasPerMission", zap.Error(err))
		return
	}
	if !has {
		return
	}

	if isAdmin && ledgerPermission.OprUsers == "" {
		if err = initAdminOprUser(
			botModel.BotID,
			chatGroupID,
			userID,
		); err != nil {
			global.GVA_LOG.Error("initAdminOprUser", zap.Error(err))
			return
		}
		if err = ledgerPermission.Release(); err != nil {
			global.GVA_LOG.Warn("Ledger CacheDelItem failed", zap.Error(err))
		}
		permit = true
	}

	if !ledgerPermission.HasUserPermission(userID, userName) {
		return
	}

	permit = true
	return
}

// 是否是admin 用户
func IsChatAdmin(botToken string, chatID int64, userID int64) (bool, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return false, err
	}

	cfg := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	}

	member, err := bot.GetChatMember(cfg)
	if err != nil {
		return false, err
	}

	switch member.Status {
	case "creator", "administrator":
		return true, nil
	default:
		return false, nil
	}
}

func initAdminOprUser(
	botID int64,
	chatGroupID int64,
	userID int64,
) error {
	return global.GVA_MYSQL.Model(&ledger.LedgerPermission{}).
		Where("bot_id = ? AND chat_group_id = ?", botID, chatGroupID).
		Where("(opr_users IS NULL OR opr_users = '')").
		Update("opr_users", strconv.FormatInt(userID, 10)).
		Error
}
