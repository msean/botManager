package handle

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	"go.uber.org/zap"
)

func SyncChatGroup(botModel bot.Bot, tgMsg tgbotapi.Update) {
	if tgMsg.Message == nil {
		return
	}
	chatID := tgMsg.Message.Chat.ID
	chatName := tgMsg.Message.Chat.Title

	if chatID == 0 || chatName == "" {
		return
	}

	cacheObject := cache.NewBotChatGroupCache(bot.BotChatGroup{
		BotID:       botModel.BotID,
		ChatGroupID: chatID,
	})
	has, err := cache.CacheGetItem(cacheObject)
	if err != nil {
		global.GVA_LOG.Error("CacheGet failed", zap.Int64("chatID", chatID), zap.Error(err))
		return
	}

	// 如果缓存或数据库不存在，创建新记录
	if !has {
		newGroup := bot.BotChatGroup{
			BotID:         botModel.BotID,
			ChatGroupID:   chatID,
			ChatGroupName: chatName,
		}
		if createErr := global.GVA_DB.Create(&newGroup).Error; createErr != nil {
			global.GVA_LOG.Error("failed to create new chat group",
				zap.Int64("chatID", chatID),
				zap.String("chatName", chatName),
				zap.Error(createErr),
			)
			return
		}
		global.GVA_LOG.Info("new chat group added",
			zap.Int64("chatID", chatID),
			zap.String("chatName", chatName),
		)
		return
	}

	// 如果数据库/缓存存在，但名称不同，则更新数据库和缓存
	if cacheObject.ChatGroupName != chatName {
		if err := global.GVA_DB.Model(&bot.BotChatGroup{}).
			Where("bot_id = ? AND chat_group_id = ?", botModel.BotID, chatID).
			Update("chat_group_name", chatName).Error; err != nil {
			global.GVA_LOG.Error("failed to update chat group name",
				zap.Int64("chatID", chatID),
				zap.String("newName", chatName),
				zap.Error(err),
			)
		} else {
			global.GVA_LOG.Info("chat group name updated",
				zap.Int64("chatID", chatID),
				zap.String("newName", chatName),
			)
		}

		if err = cache.CacheDelete(cacheObject); err != nil {
			global.GVA_LOG.Error("failed to update chat group name",
				zap.Int64("chatID", chatID),
				zap.String("newName", chatName),
				zap.Error(err),
			)
		}
	}
}

func SyncChannel(botModel bot.Bot, tgMsg tgbotapi.Update) {
	msg := tgMsg.ChannelPost

	if msg == nil || msg.Chat == nil {
		return
	}

	channelChatID := msg.Chat.ID
	channelChatName := msg.Chat.Title

	global.GVA_LOG.Debug("SyncChannelMsg",
		zap.Int64("chatID", channelChatID),
	)

	cacheObject := cache.NewBotChannelCache(bot.BotChannel{
		BotID:     botModel.BotID,
		ChannelID: channelChatID,
	})
	has, err := cache.CacheGetItem(cacheObject)
	if err != nil {
		global.GVA_LOG.Error("CacheGet failed", zap.Int64("chatID", channelChatID), zap.Error(err))
		return
	}

	// 如果缓存或数据库不存在，创建新记录
	if !has {
		newGroup := bot.BotChannel{
			BotID:       botModel.BotID,
			ChannelID:   channelChatID,
			ChannelName: channelChatName,
		}
		if createErr := global.GVA_DB.Create(&newGroup).Error; createErr != nil {
			global.GVA_LOG.Error("failed to create new chat group",
				zap.Int64("chatID", channelChatID),
				zap.String("chatName", channelChatName),
				zap.Error(createErr),
			)
			return
		}
		global.GVA_LOG.Info("new chat group added",
			zap.Int64("chatID", channelChatID),
			zap.String("chatName", channelChatName),
		)
		return
	}

	// 如果数据库/缓存存在，但名称不同，则更新数据库和缓存
	if cacheObject.ChannelName != channelChatName {
		if err := global.GVA_DB.Model(&bot.BotChannel{}).
			Where("bot_id = ? AND channel_id = ?", botModel.BotID, channelChatID).
			Update("channel_name", channelChatName).Error; err != nil {
			global.GVA_LOG.Error("failed to update chat group name",
				zap.Int64("chatID", channelChatID),
				zap.String("newName", channelChatName),
				zap.Error(err),
			)
		} else {
			global.GVA_LOG.Info("chat group name updated",
				zap.Int64("chatID", channelChatID),
				zap.String("newName", channelChatName),
			)
		}

		if err = cacheObject.Release(); err != nil {
			global.GVA_LOG.Error("failed to update chat group name",
				zap.Int64("chatID", channelChatID),
				zap.String("newName", channelChatName),
				zap.Error(err),
			)
		}
	}
}
