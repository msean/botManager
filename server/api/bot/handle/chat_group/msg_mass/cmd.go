package msgmass

import (
	"errors"
	"strconv"
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	Parser interface {
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
		parsers []Parser
	}
	ShouldPermissionAwareWithOutAdmin struct{}
	NotPermissionAwareWithAdmin       struct{}
)

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

func NewParserChain(parsers ...Parser) *ParserChain {
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
		&BotAtUsersSetHandler{},
		&BotAtUserShowHandler{},
	)

	return chain.Handle(botModel, tgMsg)
}

// 检测非admin用户是否有权限
func HasPerMission(botModel bot.Bot, tgMsg botapi.Update) (
	has bool,
	permit bool,
	err error,
) {
	chatGroupID := tgMsg.Message.Chat.ID
	userID := tgMsg.Message.From.ID
	userName := tgMsg.Message.From.UserName

	var permission bot.BotMassMsgPermission

	err = global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		botModel.BotID,
		chatGroupID,
	).First(&permission).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没配置权限
			return false, false, nil
		}
		return false, false, err
	}

	has = true

	// 空表示没人有权限
	if strings.TrimSpace(permission.OprUsers) == "" {
		return has, false, nil
	}

	// 拆分 opr_users
	users := strings.Split(permission.OprUsers, ",")

	for _, u := range users {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}

		// 1️⃣ 匹配 userID
		if strconv.FormatInt(userID, 10) == u {
			return has, true, nil
		}

		// 2️⃣ 匹配 username
		if userName != "" && u == userName {
			return has, true, nil
		}
	}

	return has, false, nil
}

// 是否是admin 用户
func IsChatAdmin(botToken string, chatID int64, userID int64, userName string) (bool, error) {
	// 写死
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
