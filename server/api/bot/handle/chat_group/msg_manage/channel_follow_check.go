package msgmanage

import (
	"strconv"
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func ChannelFollowCheck(
	botModel bot.Bot,
	chatGroup bot.BotChatGroup,
	update *botapi.Update,
) {
	if chatGroup.MustJoinChannels == "" || chatGroup.InvaidChannelFoldLink == "" {
		return
	}

	if update.Message == nil || update.Message.From == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// ① DB 缓存
	var record bot.BotChatGroupRelatedChannelFollow
	err := global.GVA_MYSQL.
		Where("bot_id = ? AND chat_group_id = ? AND user_id = ?",
			botModel.BotID,
			chatGroup.ChatGroupID,
			userID,
		).
		First(&record).Error

	if err == nil {
		return
	}

	botAPI, err := bot_handler.NewBot(botModel.Token)
	if err != nil {
		return
	}

	channelIDs := strings.Split(chatGroup.MustJoinChannels, ",")

	for _, chIDStr := range channelIDs {
		chID, _ := strconv.ParseInt(strings.TrimSpace(chIDStr), 10, 64)

		ok, _ := IsUserJoinedChannel(botAPI, chID, userID)
		if !ok {
			// ❌ 禁言 10 分钟
			_ = BanUser(
				botModel,
				*update,
				constant.BanTypeUnFollowChannel,
				10*time.Minute,
			)

			sendMustJoinMessage(
				botAPI,
				chatID,
				update.Message.MessageID,
				chatGroup.ChatGroupID,
				chatGroup.InvaidChannelFoldLink,
			)
			return
		}
	}

	// ✅ 全通过，写 DB
	_ = global.GVA_MYSQL.Create(&bot.BotChatGroupRelatedChannelFollow{
		UserID:      userID,
		BotID:       botModel.BotID,
		ChatGroupID: chatGroup.ChatGroupID,
		CheckTime:   time.Now(),
	}).Error
}
func sendMustJoinMessage(
	botAPI *bot_handler.Bot,
	chatID int64,
	replyMsgID int,
	chatGroupID int64,
	channelFoldLink string,
) {
	msg := botapi.NewMessage(
		chatID,
		"🚫 若您要在群中发言，请先关注我们的频道\n\n"+
			"👉 完成后点击「我已订阅」即可解除禁言",
	)

	msg.ReplyToMessageID = replyMsgID

	msg.ReplyMarkup = botapi.NewInlineKeyboardMarkup(
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonURL(
				"🔔 一键关注频道",
				channelFoldLink,
			),
		),
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				"✅ 我已订阅",
				"check_join:"+strconv.FormatInt(chatGroupID, 10),
			),
		),
	)

	_, _ = botAPI.Send(msg)
}

func IsUserJoinedChannel(
	bot *bot_handler.Bot,
	channelID int64,
	userID int64,
) (bool, error) {

	cfg := botapi.GetChatMemberConfig{
		ChatConfigWithUser: botapi.ChatConfigWithUser{
			ChatID: channelID,
			UserID: userID,
		},
	}

	member, err := bot.GetChatMember(cfg)
	if err != nil {
		return false, err
	}

	switch member.Status {
	case "creator", "administrator", "member", "restricted":
		return true, nil
	case "left", "kicked":
		return false, nil
	default:
		// 兜底，未来 Telegram 加状态也不炸
		return false, nil
	}
}

func handleCallback(
	botModel bot.Bot,
	update *botapi.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	if !strings.HasPrefix(data, "check_join:") {
		return
	}

	chatGroupID, err := strconv.ParseInt(
		strings.TrimPrefix(data, "check_join:"),
		10,
		64,
	)
	if err != nil {
		return
	}

	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Chat.ID

	var chatGroup bot.BotChatGroup
	if err := global.GVA_MYSQL.
		Where("chat_group_id = ?", chatGroupID).
		First(&chatGroup).Error; err != nil {
		return
	}

	botAPI, err := bot_handler.NewBot(botModel.Token)
	if err != nil {
		return
	}

	channelIDs := strings.Split(chatGroup.MustJoinChannels, ",")

	for _, chIDStr := range channelIDs {
		chID, _ := strconv.ParseInt(strings.TrimSpace(chIDStr), 10, 64)

		ok, _ := IsUserJoinedChannel(botAPI, chID, userID)
		if !ok {
			msg := botapi.NewMessage(
				chatID,
				"❌ 仍有频道未订阅，请先完成关注后再点击「我已订阅」",
			)
			msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
			_, _ = botAPI.Send(msg)
			return
		}
	}

	// ✅ 解禁

	// 写 DB 缓存
	_ = global.GVA_MYSQL.Create(&bot.BotChatGroupRelatedChannelFollow{
		UserID:      userID,
		BotID:       botModel.BotID,
		ChatGroupID: chatGroupID,
		CheckTime:   time.Now(),
	}).Error

	msg := botapi.NewMessage(
		chatID,
		"✅ 验证通过，已解除禁言，现在可以正常发言了",
	)
	msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	_, _ = botAPI.Send(msg)
}
