package channel

import (
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func Entrance(botModel bot.Bot, tgMsg botapi.Update) (err error) {
	SyncChannel(botModel, tgMsg)
	return
}

func SyncChannel(botModel bot.Bot, tgMsg botapi.Update) {
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
