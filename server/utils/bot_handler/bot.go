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

	cfg := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       false,
			CanSendMediaMessages:  false,
			CanSendPolls:          false,
			CanSendOtherMessages:  false,
			CanAddWebPagePreviews: false,
			CanChangeInfo:         false,
			CanInviteUsers:        false,
			CanPinMessages:        false,
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

	cfg := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true, // true 表示丢弃所有未处理的消息
	}

	if _, err = bot.Request(cfg); err != nil {
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

func UnRegisterWebhook(botToken string, dropPending bool) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	cfg := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: dropPending, // true 表示丢弃所有未处理的消息
	}

	_, err = bot.Request(cfg)
	return err
}

func (b *Bot) DeleteMsg(chatID int64, msgID int) (err error) {
	botAPI, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return err
	}

	cfg := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: msgID,
	}
	_, err = botAPI.Request(cfg)
	return
}

func SendTextMessage(chatID int64, token string, text string) {
	bot, _ := tgbotapi.NewBotAPI(token) // 你可以传 token
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
