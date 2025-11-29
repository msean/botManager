package bot_handler

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	token  string
	botApi *tgbotapi.BotAPI
}

func NewBot(token string) (bot *Bot, err error) {
	bot = &Bot{
		token: token,
	}
	if bot.botApi, err = tgbotapi.NewBotAPI(token); err != nil {
		return
	}
	return
}

func (b *Bot) BanUser(chatID, userID int64, until int64) (err error) {
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
	_, err = b.botApi.Request(cfg)
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
	cfg := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: msgID,
	}
	_, err = b.botApi.Request(cfg)
	return
}

func (b *Bot) SendTextMessage(chatID int64, token string, text string) (err error) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err = b.botApi.Send(msg)
	return
}

func (b *Bot) SendMarkDownMessage(chatID int64, token string, text string) (err error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	_, err = b.botApi.Send(msg)
	return
}

func EscapeMarkdownV2CodeBlock(text string) string {
	specialChars := []string{"`", "\\"}
	for _, ch := range specialChars {
		text = strings.ReplaceAll(text, ch, "\\"+ch)
	}
	return text
}

// EscapeMarkdownV2 用于普通文本，转义 MarkdownV2 保留字符
func EscapeMarkdownV2(text string) string {
	specialChars := []string{
		"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}",
		".", "!", // 普通文本里 . 和 ! 需要转义
	}
	for _, ch := range specialChars {
		text = strings.ReplaceAll(text, ch, "\\"+ch)
	}
	return text
}

// FormatRechargeMessage 生成安全的充值提示 MarkdownV2 消息
func FormatRechargeMessage(orderID uint, amount, paymentAddr string, createdAt string, leftPaidMinutes int) string {
	return fmt.Sprintf(
		"订单号：%d\n"+
			"转账金额：`%s` USDT （点击即可复制）\n"+
			"转账地址：`%s` （点击即可复制）\n"+
			"充值时间：%s\n\n"+
			"⚠️注意：\n"+
			"▫️注意小数点 %s 转错金额不能到账\n"+
			"▫️请在%d分钟完成付款，转错金额不能到账。\n\n"+
			"转账%d分钟后没到账及时联系",
		orderID,
		EscapeMarkdownV2CodeBlock(amount),
		EscapeMarkdownV2CodeBlock(paymentAddr),
		createdAt,
		EscapeMarkdownV2(amount),
		leftPaidMinutes,
		leftPaidMinutes,
	)
}

// SendAdMessage 统一发送广告内容（文字 / 图片 + 文本 / 视频 + 文本）
// 如果 replyMarkup != nil，则用作按钮，否则不带按钮
func (b *Bot) TgSend(token string, chatID int64, medias []MediaItem, replyMarkup interface{}) (tgMsg tgbotapi.Message, err error) {
	var caption string
	var photoID string
	var videoID string

	// 整理用户发送的数据
	for _, m := range medias {
		switch m.Type {
		case "text":
			caption = m.Text
		case "photo":
			photoID = m.FileID
		case "video":
			videoID = m.FileID
		}
	}

	// PHOTO
	if photoID != "" {
		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(photoID))
		msg.Caption = caption

		if replyMarkup != nil {
			msg.ReplyMarkup = replyMarkup
		}

		return b.botApi.Send(msg)
	}

	// VIDEO
	if videoID != "" {
		msg := tgbotapi.NewVideo(chatID, tgbotapi.FileID(videoID))
		msg.Caption = caption

		if replyMarkup != nil {
			msg.ReplyMarkup = replyMarkup
		}

		return b.botApi.Send(msg)
	}

	// ONLY TEXT
	msg := tgbotapi.NewMessage(chatID, caption)

	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}

	return b.botApi.Send(msg)
}
