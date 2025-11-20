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
			tgbotapi.NewKeyboardButton("⚡能源租用"),
			tgbotapi.NewKeyboardButton("💵能量理财"),
			tgbotapi.NewKeyboardButton("📅888租号"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔥智能托管"),
			tgbotapi.NewKeyboardButton("✏️笔数套餐"),
			tgbotapi.NewKeyboardButton("💎会员星星"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅TRX闪兑"),
			tgbotapi.NewKeyboardButton("🛠常用功能"),
			tgbotapi.NewKeyboardButton("📣供需广告"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🦊实用导航"),
			tgbotapi.NewKeyboardButton("🏦钱包功能"),
			tgbotapi.NewKeyboardButton("🌐新浪头条"),
		),
	)

	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "欢迎使用，请选择功能：")
	msg.ReplyMarkup = keyboard

	botAPI.Send(msg)
}
