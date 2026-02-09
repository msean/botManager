package msgmass

import (
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

type (
	MessageParser interface {
		Match(botModel bot.Bot, update botapi.Update) bool
		Handle() error
	}
	PermissionAware interface {
		NeedAdmin() bool
	}
)

type (
	ParserChain struct {
		parsers []MessageParser
	}
	NeedAdminAware struct{}
)

func (c *NeedAdminAware) NeedAdmin() bool {
	return true
}

func (c *ParserChain) Handle(botModel bot.Bot, update botapi.Update) error {

	chatGroupID := update.Message.Chat.ID

	for _, parser := range c.parsers {
		global.GVA_LOG.Debug("msgmass ParserChain Handle4", zap.Any("bot", botModel))
		if !parser.Match(botModel, update) {
			continue
		}

		global.GVA_LOG.Debug("msgmass ParserChain Handle5", zap.Any("parser", parser))
		// 权限感知
		if p, ok := parser.(PermissionAware); ok {

			global.GVA_LOG.Debug("msgmass ParserChain Handle6", zap.Any("parser", parser))
			var isAdminUser bool
			isAdminUser, err := IsChatAdmin(
				botModel.Token,
				chatGroupID,
				update.Message.From.ID,
				update.Message.From.UserName,
			)
			global.GVA_LOG.Debug("msgmass ParserChain Handle7", zap.Any("err", err), zap.Any("ok", ok))
			if err != nil || !ok {
				return nil // 静默失败，不给回复
			}

			global.GVA_LOG.Debug("msgmass ParserChain Handle8", zap.Any("isAdminUser", isAdminUser))
			// 管理员权限（最高优先级）
			if p.NeedAdmin() && !isAdminUser {
				global.GVA_LOG.Info("msgmass ParserChain not isAdminUser", zap.Any("user", update.Message.Chat.UserName))
				return nil
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
		&BotAtUsersSetHandler{},
		&BotAtUserShowHandler{},
	)

	return chain.Handle(botModel, tgMsg)
}

// 是否是admin 用户
func IsChatAdmin(botToken string, chatID int64, userID int64, userName string) (bool, error) {
	// // 写死
	// if userID == 7449031746 || userID == 8099503790 || userName == "xmpaymo" {
	// 	return true, nil
	// }

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
