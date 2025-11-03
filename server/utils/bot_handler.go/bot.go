package bot_handler

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	token string
}

func NewBot(token string) *Bot {
	return &Bot{
		token: token,
	}
}

func (b *Bot) BanUser(chatID, userID int64, duration time.Duration) error {
	botAPI, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return err
	}

	until := time.Now().Add(duration).Unix()

	cfg := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: until,
	}

	_, err = botAPI.Request(cfg)
	return err
}
