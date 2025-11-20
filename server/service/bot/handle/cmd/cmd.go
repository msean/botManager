package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	TgMsgHandler *Handler
)

type (
	HandlerFunc func(update tgbotapi.Update, token string, botID int64)
	Handler     struct {
		routes map[string]HandlerFunc
	}
)

func NewHandler() *Handler {
	return &Handler{
		routes: make(map[string]HandlerFunc),
	}
}

func (r *Handler) Register(cmd string, handler HandlerFunc) {
	r.routes[cmd] = handler
}

func (r *Handler) Handle(update tgbotapi.Update, token string, botID int64) {

	if update.Message == nil {
		return
	}

	text := update.Message.Text

	if handler, ok := r.routes[text]; ok {
		handler(update, token, botID)
		return
	}
}

func InitTgMsgHandler() {
	TgMsgHandler = NewHandler()
	TgMsgHandler.Register("/start", StartHandlerfunc)
}
