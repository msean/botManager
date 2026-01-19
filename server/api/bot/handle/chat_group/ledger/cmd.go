package ledger

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
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
		SetPermission(*cache.LedgerPermissionCache)
	}
)

type (
	ParserChain struct {
		parsers []MessageParser
	}
	ShouldPermissionAware struct {
		confModel *cache.LedgerPermissionCache
	}
	NotPermissionAware struct {
		confModel *cache.LedgerPermissionCache
	}
)

func (l *ShouldPermissionAware) NeedPermission() bool {
	return true
}

func (l *ShouldPermissionAware) SetPermission(p *cache.LedgerPermissionCache) {
	l.confModel = p
}

func (l *NotPermissionAware) NeedPermission() bool {
	return false
}

func (l *NotPermissionAware) SetPermission(p *cache.LedgerPermissionCache) {
	l.confModel = p
}

func (c *ParserChain) Handle(botModel bot.Bot, update tgbotapi.Update) error {

	for _, parser := range c.parsers {
		if !parser.Match(botModel, update) {
			continue
		}

		// 是否需要权限
		if p, ok := parser.(PermissionAware); ok && p.NeedPermission() {
			permission, permit, err := HasPerMission(botModel, update)
			if err != nil || !permit {
				return err
			}
			p.SetPermission(permission)
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

func HasPerMission(botModel bot.Bot, tgMsg tgbotapi.Update) (ledgerPermission *cache.LedgerPermissionCache, permit bool, err error) {

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
	if !ledgerPermission.HasUserPermission(userID, userName) {
		return
	}

	permit = true
	return
}
