package cmd

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/service/cache"
)

func BuildReplyKeyboard(cmdButtons json.RawMessage) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]struct {
		Name    string `json:"name"`
		BindCmd string `json:"bindCmd"`
	}

	json.Unmarshal(cmdButtons, &rows)

	keyboardRows := make([][]tgbotapi.KeyboardButton, 0)

	for _, row := range rows {
		btnRow := make([]tgbotapi.KeyboardButton, 0)
		for _, b := range row {
			btnRow = append(btnRow, tgbotapi.NewKeyboardButton(b.BindCmd))
			// 或显示名字 btnRow = append(btnRow, tgbotapi.NewKeyboardButton(b.Name))
		}
		keyboardRows = append(keyboardRows, btnRow)
	}

	return tgbotapi.NewReplyKeyboard(keyboardRows...)
}

func StartHandlerfunc(update tgbotapi.Update, token string, cfg cache.BotCmdCache) {
	bot, _ := tgbotapi.NewBotAPI(token)
	chatID := update.Message.Chat.ID

	keyboard := BuildReplyKeyboard(cfg.CmdButtons)

	// 3. 回复欢迎语 + 键盘
	msg := tgbotapi.NewMessage(chatID, cfg.Content)
	msg.ReplyMarkup = keyboard

	bot.Send(msg)
}
