package bot_handler

import (
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

func (b *Bot) BanUser(chatID, userID int64, until int64) error {
	botAPI, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return err
	}

	// until := time.Now().Add(duration).Unix()
	// global.GVA_LOG.Info("util", zap.Int("util", int(until)))

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

func RegisterWebhook(botToken, webhookURL string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return err
	}

	_, err = bot.Request(wh)
	if err != nil {
		return err
	}
	return nil
}
