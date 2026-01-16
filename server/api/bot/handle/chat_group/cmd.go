package chat_group

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageParser interface {
	Match(botID int64, update tgbotapi.Update) bool
	Handle() error
}

type ParserChain struct {
	parsers []MessageParser
}

func (c *ParserChain) Handle(botID int64, update tgbotapi.Update) error {
	for _, parser := range c.parsers {
		if parser.Match(botID, update) {
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

func Handle(botID int64, tgMsg tgbotapi.Update) (err error) {
	chain := NewParserChain(
		&Ledger{},
	)

	return chain.Handle(botID, tgMsg)
}
