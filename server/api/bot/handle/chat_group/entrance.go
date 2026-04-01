package chatgroup

import (
	"github.com/msean/botmanager/server/api/bot/handle/chat_group/ledger"
	"github.com/msean/botmanager/server/api/bot/handle/chat_group/ledger2"
	msgmanage "github.com/msean/botmanager/server/api/bot/handle/chat_group/msg_manage"
	msgmass "github.com/msean/botmanager/server/api/bot/handle/chat_group/msg_mass"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botSvc "github.com/msean/botmanager/server/service/bot"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func Entrance(botModel bot.Bot, tgMsg botapi.Update) (err error) {
	global.GVA_LOG.Debug("CacheGet failed", zap.Any("redis", global.GVA_REDIS))
	var has bool
	chatGroupID := bot_handler.GetChatID(tgMsg)
	chatGroup := cache.NewBotChatGroupCache(botModel.BotID, chatGroupID)
	has, err = cache.CacheGetItem(chatGroup)
	if err != nil {
		global.GVA_LOG.Error("CacheGet failed", zap.Int64("chatID", chatGroupID), zap.Error(err))
		return
	}

	global.GVA_LOG.Debug("HandelChatGroup Sync", zap.Int64("SyncMessage", chatGroup.SyncMessage))
	global.GVA_LOG.Debug("BotHandler HandelChatGroup", zap.Any("bot", botModel))

	// 同步群组信息
	SyncChatGroup(botModel, tgMsg, chatGroup, has)

	// 需要同步群聊消息
	go func() {
		if chatGroup.SyncMessage != 1 {
			return
		}
		defer func() {
			if r := recover(); r != nil {
				global.GVA_LOG.Error(
					"panic in SyncChatGroupMessage",
					zap.Any("recover", r),
					zap.Int64("chatGroupID", chatGroupID),
					zap.Stack("stack"),
				)
			}
		}()

		SyncChatGroupMessage(botModel.BotID, chatGroupID, tgMsg)
	}()

	// 记账功能入口
	go func() {
		if botModel.IsForLedger != 1 {
			return
		}
		global.GVA_LOG.Debug("in IsForLedger")
		defer func() {
			if r := recover(); r != nil {
				global.GVA_LOG.Error(
					"panic in chat_group Handle",
					zap.Any("recover", r),
					zap.Int64("chatGroupID", chatGroupID),
					zap.Stack("stack"),
				)
			}
		}()

		ledger.Handle(botModel, tgMsg)
	}()

	// 消息管理入口
	go func() {
		if botModel.IsForMsgMgr != 1 {
			return
		}
		global.GVA_LOG.Debug("in msgmanage")
		defer func() {
			if r := recover(); r != nil {
				global.GVA_LOG.Error(
					"panic in ban msg Handle",
					zap.Any("recover", r),
					zap.Int64("chatGroupID", chatGroupID),
					zap.Stack("stack"),
				)
			}
		}()

		msgmanage.Dispatch(botModel, tgMsg, *chatGroup)
	}()

	// 消息管理入口
	go func() {
		if botModel.IsForMsgMgr != 1 {
			return
		}
		global.GVA_LOG.Debug("in msgmass")
		defer func() {
			if r := recover(); r != nil {
				global.GVA_LOG.Error(
					"panic in ban msg Handle",
					zap.Any("recover", r),
					zap.Int64("chatGroupID", chatGroupID),
					zap.Stack("stack"),
				)
			}
		}()

		msgmass.Handle(botModel, tgMsg)
	}()

	// 消息管理入口
	go func() {
		if botModel.IsForLedger2 != 1 {
			return
		}
		global.GVA_LOG.Debug("in IsForLedger2")
		defer func() {
			if r := recover(); r != nil {
				global.GVA_LOG.Error(
					"panic in ban msg Handle",
					zap.Any("recover", r),
					zap.Int64("chatGroupID", chatGroupID),
					zap.Stack("stack"),
				)
			}
		}()

		ledger2.Dispatch(botModel, tgMsg, *chatGroup)
	}()

	return
}

func SyncChatGroup(botModel bot.Bot, tgMsg botapi.Update, chatGroup *cache.BotChatGroupCache, has bool) (err error) {
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
			SyncMessage:   1,  // 默认不开启吧
			MaxWords:      -1, // 无限制
			BanForward:    1,  // 禁用
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

// 同步群消息
func SyncChatGroupMessage(botID int64, chatGroupID int64, tgMsg botapi.Update) (err error) {
	svc := botSvc.NewBotChatHistorySvc(botID, chatGroupID)
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
