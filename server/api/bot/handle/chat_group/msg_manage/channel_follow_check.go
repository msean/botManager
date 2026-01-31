package msgmanage

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
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

	global.GVA_LOG.Debug("ChannelFollowCheck11", zap.Any("botModel", botModel), zap.Error(err))
	if err == nil {
		global.GVA_LOG.Error("ChannelFollowCheck", zap.Any("botModel", botModel), zap.Error(err))
		return
	}

	botAPI, err := bot_handler.NewBot(botModel.Token)
	if err != nil {
		global.GVA_LOG.Error("ChannelFollowCheck NewBot", zap.Any("botModel", botModel), zap.Error(err))
		return
	}

	if !shouldSkipChannelCheck(botAPI, chatID, userID) {
		return
	}

	channelIDs := strings.Split(chatGroup.MustJoinChannels, ",")
	global.GVA_LOG.Debug("ChannelFollowCheck22", zap.Any("botModel", botModel), zap.Error(err))

	for index, chIDStr := range channelIDs {
		chID, _ := strconv.ParseInt(strings.TrimSpace(chIDStr), 10, 64)

		ok, getErr := IsUserJoinedChannel(botAPI, chID, userID)

		// 防止429 并发请求错误
		if index != len(channelIDs)-1 {
			time.Sleep(100 * time.Millisecond)
		}

		global.GVA_LOG.Debug(
			"ChannelFollowCheck IsUserJoinedChannel",
			zap.Int64("channelID", chID),
			zap.Int64("userID", userID),
			zap.Bool("ok", ok),
			zap.Error(getErr),
		)

		// 🚨 ① 只要有一次请求异常，直接放弃本轮校验
		if getErr != nil {
			global.GVA_LOG.Warn(
				"ChannelFollowCheck abort due to getChatMember error",
				zap.Int64("channelID", chID),
				zap.Int64("userID", userID),
				zap.Error(getErr),
			)
			return
		}

		if !ok {
			sendMustJoinMessage(
				botAPI,
				chatID,
				userID,
				update.Message.MessageID,
				chatGroup.ChatGroupID,
				chatGroup.InvaidChannelFoldLink,
			)
			BanUser(
				botModel,
				*update,
				constant.BanTypeUnFollowChannel,
				10*time.Minute,
			)
			return
		}

		// ✅ ③ 已关注 → 继续检查下一个
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
	userID int64,
	replyMsgID int,
	chatGroupID int64,
	channelFoldLink string,
) (err error) {
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
				fmt.Sprintf("check_join:%d:%d", chatGroupID, userID),
			),
		),
	)

	global.GVA_LOG.Debug("ChannelFollowCheck sendMustJoinMessage",
		zap.Any("chatGroupID", chatGroupID),
		zap.Any("channelFoldLink", channelFoldLink),
		zap.Any("chatID", chatID))

	if _, err = botAPI.Send(msg); err != nil {
		global.GVA_LOG.Error("ChannelFollowCheck sendMustJoinMessage",
			zap.Any("chatGroupID", chatGroupID),
			zap.Any("channelFoldLink", channelFoldLink),
			zap.Any("chatID", chatID),
			zap.Error(err))
	}
	return
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
		global.GVA_LOG.Error("ChannelFollowCheck IsUserJoinedChannel", zap.Any("channelID", channelID), zap.Any("", userID), zap.Error(err))
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

	splits := strings.Split(data, ":")
	if len(splits) != 3 {
		return
	}
	chatGroupID, err := strconv.ParseInt(
		splits[1],
		10,
		64,
	)
	if err != nil {
		return
	}

	// _userID, err := strconv.ParseInt(
	// 	splits[2],
	// 	10,
	// 	64,
	// )
	// if err != nil {
	// 	return
	// }

	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Chat.ID

	// if userID != _userID {
	// 	return
	// }

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

	global.GVA_LOG.Error("ChannelFollowCheck IsUserJoinedChannel", zap.Any("Env", global.GVA_CONFIG.System.Env), zap.Any("", userID), zap.Error(err))
	if global.GVA_CONFIG.System.Env != "test" {
		// 写 DB 缓存
		_ = global.GVA_MYSQL.Create(&bot.BotChatGroupRelatedChannelFollow{
			UserID:      userID,
			BotID:       botModel.BotID,
			ChatGroupID: chatGroupID,
			CheckTime:   time.Now(),
		}).Error
	}
	if err := botAPI.UnMuteUser(chatID, userID); err == nil {
		msg := botapi.NewMessage(
			chatID,
			"✅ 验证通过，已解除禁言，现在可以正常发言了",
		)
		msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
		_, _ = botAPI.Send(msg)
	}
}

func shouldSkipChannelCheck(
	botAPI *bot_handler.Bot,
	chatID int64,
	userID int64,
) bool {

	cfg := botapi.GetChatMemberConfig{
		ChatConfigWithUser: botapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	}

	member, err := botAPI.GetChatMember(cfg)
	if err != nil {
		// ⚠️ 出错不跳过，避免被 API 异常绕过校验
		global.GVA_LOG.Warn(
			"shouldSkipChannelCheck GetChatMember error",
			zap.Int64("chatID", chatID),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return false
	}

	// ① 必须是 admin / creator
	if member.Status != "administrator" && member.Status != "creator" {
		return true
	}

	// ② 必须能拿到入群时间
	// if member.JoinedAt == 0 {
	// 	return false
	// }

	// joinedAt := time.Unix(int64(member.JoinedAt), 0)

	// // ③ 入群 ≥ 2 个月
	// if time.Since(joinedAt) >= 60*24*time.Hour {
	// 	return true
	// }

	return false
}
