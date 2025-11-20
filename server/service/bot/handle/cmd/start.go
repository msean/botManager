package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"go.uber.org/zap"
)

func StartHandlerfunc(update tgbotapi.Update, token string, botID int64) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		global.GVA_LOG.Error("StartHandlerfunc", zap.Error(err))
		return
	}
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚡充值"),
			tgbotapi.NewKeyboardButton("💵充值价格"),
		),
	)

	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "欢迎使用，请选择功能：")
	msg.ReplyMarkup = keyboard

	botAPI.Send(msg)
}
