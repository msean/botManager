package ledger

import (
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type (
	MessageParser interface {
		Match(botModel bot.Bot, update botapi.Update) bool
		Handle() error
	}
	ParserChain struct {
		parsers []MessageParser
	}
)

func (c *ParserChain) Handle(botModel bot.Bot, update botapi.Update) error {

	for _, parser := range c.parsers {
		if !parser.Match(botModel, update) {
			continue
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

// 解析
func Handle(botModel bot.Bot, tgMsg botapi.Update) (err error) {
	if tgMsg.Message == nil || tgMsg.Message.Text == "" {
		return
	}
	chain := NewParserChain(
		&TronAddressParser{},
		&GroupListParser{},
		&GroupStatParser{},
	)
	return chain.Handle(botModel, tgMsg)
}
