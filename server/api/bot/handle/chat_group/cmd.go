package chat_group

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/model/bot"
)

type MessageParser interface {
	Match(botModel bot.Bot, update tgbotapi.Update) bool
	Handle() error
}

type ParserChain struct {
	parsers []MessageParser
}

func (c *ParserChain) Handle(botModel bot.Bot, update tgbotapi.Update) error {
	for _, parser := range c.parsers {
		if parser.Match(botModel, update) {
			return parser.Handle()
		}
	}
	return nil
}

func NewParserChain(parsers ...MessageParser) *ParserChain {
	return &ParserChain{
		parsers: parsers,
	}
}

func Handle(botModel bot.Bot, tgMsg tgbotapi.Update) (err error) {
	chain := NewParserChain(
		&LedgerHandler{},
		&LedgerShowFeeRateHandler{},
		&LedgerSetFeeRateHandler{},
	)

	return chain.Handle(botModel, tgMsg)
}
