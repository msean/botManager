package handle

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botSvc "github.com/msean/botmanager/server/service/bot"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

func SyncChatGroup(botModel bot.Bot, tgMsg tgbotapi.Update, chatGroup *cache.BotChatGroupCache, has bool) (err error) {
	if tgMsg.Message == nil {
		return
	}
	chatID := tgMsg.Message.Chat.ID
	chatName := tgMsg.Message.Chat.Title

	if chatID == 0 || chatName == "" {
		return
	}

	// 如果没有传 则去获取
	if chatGroup == nil {
		chatGroupID := bot_handler.GetChatID(tgMsg)
		chatGroup := cache.NewBotChatGroupCache(botModel.BotID, chatGroupID)
		has, err = cache.CacheGetItem(chatGroup)
		if err != nil {
			global.GVA_LOG.Error("CacheGet failed", zap.Int64("chatID", chatGroupID), zap.Error(err))
			return
		}
	}

	// 如果缓存或数据库不存在，创建新记录
	if !has {
		newGroup := bot.BotChatGroup{
			BotID:         botModel.BotID,
			ChatGroupID:   chatID,
			ChatGroupName: chatName,
		}
		if createErr := global.GVA_MYSQL.Create(&newGroup).Error; createErr != nil {
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
	if chatGroup.ChatGroupName != chatName {
		if err := global.GVA_MYSQL.Model(&bot.BotChatGroup{}).
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

		if err = cache.CacheDelete(chatGroup); err != nil {
			global.GVA_LOG.Error("failed to update chat group name",
				zap.Int64("chatID", chatID),
				zap.String("newName", chatName),
				zap.Error(err),
			)
		}
	}
	return
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

	cacheObject := cache.NewBotChannelCache(botModel.BotID, channelChatID)
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
		if createErr := global.GVA_MYSQL.Create(&newGroup).Error; createErr != nil {
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
		if err := global.GVA_MYSQL.Model(&bot.BotChannel{}).
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

func SyncChatGroupMessage(botID int64, chatGroupID int64, tgMsg tgbotapi.Update) (err error) {
	svc := botSvc.NewBotMsgRecordSvc(botID, chatGroupID)
	if err = svc.Sync(); err != nil {
		global.GVA_LOG.Error("SyncChatGroupMessage NewBotMsgRecordSvc",
			zap.Int64("chatID", chatGroupID),
			zap.Int64("botID", botID),
			zap.Any("msg", tgMsg),
			zap.Error(err),
		)
		return
	}
	if err = svc.SaveMessage(tgMsg); err != nil {
		global.GVA_LOG.Error("SyncChatGroupMessage SaveMessage",
			zap.Int64("chatID", chatGroupID),
			zap.Int64("botID", botID),
			zap.Any("msg", tgMsg),
			zap.Error(err),
		)
	}
	return
}
