package ledger

import (
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

type (
	MessageParser interface {
		Match(botModel bot.Bot, update botapi.Update) bool
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

func (c *ParserChain) Handle(botModel bot.Bot, update botapi.Update) error {

	global.GVA_LOG.Debug("ParserChain Handle", zap.Any("bot", botModel))

	global.GVA_LOG.Debug("ParserChain Handle2", zap.Any("bot", botModel))
	chatGroupID := update.Message.Chat.ID

	global.GVA_LOG.Debug("ParserChain Handle3", zap.Any("bot", botModel))

	for _, parser := range c.parsers {
		global.GVA_LOG.Debug("ParserChain Handle4", zap.Any("bot", botModel))
		if !parser.Match(botModel, update) {
			continue
		}

		global.GVA_LOG.Debug("ParserChain Handle5", zap.Any("parser", parser))
		// 权限感知
		if p, ok := parser.(PermissionAware); ok {

			global.GVA_LOG.Debug("ParserChain Handle6", zap.Any("parser", parser))
			var isAdminUser bool
			isAdminUser, err := IsChatAdmin(
				botModel.Token,
				chatGroupID,
				update.Message.From.ID,
				update.Message.From.UserName,
			)
			global.GVA_LOG.Debug("ParserChain Handle7", zap.Any("err", err), zap.Any("ok", ok))
			if err != nil || !ok {
				return nil // 静默失败，不给回复
			}

			global.GVA_LOG.Debug("ParserChain Handle8", zap.Any("isAdminUser", isAdminUser))
			// 管理员权限（最高优先级）
			if p.NeedAdmin() && !isAdminUser {
				global.GVA_LOG.Info("ParserChain not isAdminUser", zap.Any("user", update.Message.Chat.UserName))
				return nil
			}

			// 业务权限（Ledger 操作人） 非管理员用户 管理员用户可以操作任何
			if p.NeedPermission() && !isAdminUser {
				permission, has, permit, err := HasPerMission(botModel, update)
				// 没有的话直接返回
				if !has {
					return nil
				}
				if err != nil {
					return nil
				}
				if !isAdminUser && !permit {
					return nil
				}

				p.SetPermission(permission)
			}

			if isAdminUser {
				permission, has, _, err := HasPerMission(botModel, update)
				global.GVA_LOG.Debug("ParserChain Handle", zap.Any("has", has), zap.Any("err", err), zap.Any("permission", permission))
				if err != nil {
					return nil
				}

				if !has {
					// 不存在就创建
					permissionModel := ledger.LedgerPermission{
						BotID:       botModel.BotID,
						ChatGroupID: chatGroupID,
					}
					if err = global.GVA_MYSQL.Create(&permissionModel).Error; err != nil {
						return err
						// }
					}
					permission = cache.NewLedgerPermissionCache(botModel.BotID, chatGroupID)
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

func Handle(botModel bot.Bot, tgMsg botapi.Update) (err error) {
	if tgMsg.Message == nil || tgMsg.Message.Text == "" {
		return
	}
	chain := NewParserChain(
		&LedgerHandler{},
		&LedgerShowFeeRateHandler{},
		&LedgerSetFeeRateHandler{},
		&AddOprUserHandler{},
		&CleanHandler{},
	)

	return chain.Handle(botModel, tgMsg)
}

// 检测非admin用户是否有权限
func HasPerMission(botModel bot.Bot, tgMsg botapi.Update) (ledgerPermission *cache.LedgerPermissionCache, has bool, permit bool, err error) {

	chatGroupID := tgMsg.Message.Chat.ID
	userID := tgMsg.Message.From.ID
	userName := tgMsg.Message.From.UserName

	ledgerPermission = cache.NewLedgerPermissionCache(botModel.BotID, chatGroupID)

	if has, err = cache.CacheGetItem(ledgerPermission); err != nil {
		global.GVA_LOG.Error("Ledger HasPerMission", zap.Error(err))
		return
	}

	global.GVA_LOG.Debug("Ledger HasPerMission", zap.Any("has", has), zap.Error(err), zap.Any("ledgerPermission", ledgerPermission))
	if !has {
		return
	}

	if !ledgerPermission.HasUserPermission(userID, userName) {
		return
	}

	permit = true
	return
}

// 是否是admin 用户
func IsChatAdmin(botToken string, chatID int64, userID int64, userName string) (bool, error) {
	// 写死
	if userID == 7449031746 || userID == 8099503790 || userName == "xmpaymo" {
		return true, nil
	}

	bot, err := botapi.NewBotAPI(botToken)
	if err != nil {
		return false, err
	}

	cfg := botapi.GetChatMemberConfig{
		ChatConfigWithUser: botapi.ChatConfigWithUser{
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
