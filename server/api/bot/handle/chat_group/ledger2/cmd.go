package ledger2

import (
	"errors"
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	MessageParser interface {
		Match(botModel bot.Bot, update botapi.Update) bool
		Handle() error
	}
	PermissionAware interface {
		NeedPermission() bool
		NeedAdmin() bool
	}
)

type (
	ParserChain struct {
		parsers []MessageParser
	}
	ShouldPermissionAwareWithOutAdmin struct {
	}
	OnlyAdminAware struct {
	}
)

func (l *ShouldPermissionAwareWithOutAdmin) NeedPermission() bool {
	return true
}

func (l *ShouldPermissionAwareWithOutAdmin) NeedAdmin() bool {
	return false
}

func (l *OnlyAdminAware) NeedPermission() bool {
	return false
}

func (l *OnlyAdminAware) NeedAdmin() bool {
	return true
}

func (c *ParserChain) Handle(botModel bot.Bot, update botapi.Update) error {
	chatGroupID := update.Message.Chat.ID
	for _, parser := range c.parsers {
		if !parser.Match(botModel, update) {
			continue
		}
		// 权限检查
		if p, ok := parser.(PermissionAware); ok {
			var isAdminUser bool
			isAdminUser, err := IsChatAdmin(
				botModel.Token,
				chatGroupID,
				update.Message.From.ID,
				update.Message.From.UserName,
			)
			if err != nil || !ok {
				if err != nil {
					global.GVA_LOG.Error("Handle", zap.Any("update", update), zap.Error(err))
				}
				return nil
			}

			if p.NeedAdmin() && !isAdminUser {
				global.GVA_LOG.Info("ParserChain not isAdminUser", zap.Any("user", update.Message.Chat.UserName))
				return nil
			}

			// 业务权限（Ledger 操作人） 非管理员用户 管理员用户可以操作任何
			if p.NeedPermission() && !isAdminUser {
				has, permit, err := HasPerMission(botModel, update)
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
		&AddOprUserHandler{},
		&DelOprUserHandler{},
		&ListOprUserHandler{},
		&StartLedgerHandler{},
		&LedgerRecordHandler{},
		&LedgerSummaryHandler{},
	)

	return chain.Handle(botModel, tgMsg)
}

// 是否是admin 用户
func IsChatAdmin(botToken string, chatID int64, userID int64, userName string) (bool, error) {
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

func HasPerMission(botModel bot.Bot, tgMsg botapi.Update) (
	has bool,
	permit bool,
	err error,
) {
	chatGroupID := tgMsg.Message.Chat.ID
	username := tgMsg.Message.From.UserName

	var permission ledger2.LedgerPermission

	err = global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		botModel.BotID,
		chatGroupID,
	).First(&permission).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, false, nil
		}
		return false, false, err
	}

	has = true

	userList := splitUsers(permission.OprUsers)

	for _, u := range userList {
		if strings.EqualFold(u, username) {
			return true, true, nil
		}
	}

	return true, false, nil
}
