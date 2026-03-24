package handle

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/api/bot/handle/chat_group/ledger"
	"github.com/msean/botmanager/server/api/bot/handle/chat_group/ledger2"
	msgmanage "github.com/msean/botmanager/server/api/bot/handle/chat_group/msg_manage"
	msgmass "github.com/msean/botmanager/server/api/bot/handle/chat_group/msg_mass"
	"github.com/msean/botmanager/server/api/bot/handle/private"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

type BotHandler struct{}

func NewBotHandler() *BotHandler {
	return &BotHandler{}
}

func getUpdateText(u botapi.Update) string {
	if u.Message != nil {
		return u.Message.Text
	}
	if u.CallbackQuery != nil {
		return u.CallbackQuery.Data
	}
	return ""
}

func (handler *BotHandler) Handle(c *gin.Context, botID int, body []byte) (err error) {

	var tgMsg botapi.Update
	if err = json.Unmarshal(body, &tgMsg); err != nil {
		global.GVA_LOG.Error("invalid telegram tgMsg", zap.Error(err))
		return
	}

	botModel, has, err := dao.BotDao.FromBotID(global.GVA_MYSQL, botID)
	if err != nil || !has {
		global.GVA_LOG.Error("bot not found", zap.Int("botID", botID), zap.Error(err))
		return
	}

	var chatType string
	if tgMsg.Message != nil {
		chatType = tgMsg.Message.Chat.Type
	} else if tgMsg.CallbackQuery != nil {
		chatType = tgMsg.CallbackQuery.Message.Chat.Type
	} else if tgMsg.ChannelPost != nil {
		chatType = tgMsg.ChannelPost.Chat.Type
	} else {
		chatType = "unknow"
		global.GVA_LOG.Info("BotMsgHandlerSvc unkown chatType", zap.Any("tgMsg", tgMsg)) // 修复后的取文本函数
	}

	switch chatType {
	// 私聊
	case "private":
		global.GVA_LOG.Info("BotMsgHandlerSvc received msg",
			zap.Any("msg", getUpdateText(tgMsg))) // 修复后的取文本函数

		private.Handle(tgMsg, botModel.Token, int64(botModel.BotID))

	default:
		// 被拉入群
		if tgMsg.MyChatMember != nil {
			SyncChatGroup(botModel, tgMsg, nil, false)
			return nil
		}

		if tgMsg.ChannelPost != nil {
			global.GVA_LOG.Debug("BotMsgHandlerSvc  ChannelPost",
				zap.Any("msg", getUpdateText(tgMsg))) // 修复后的取文本函数
			handler.HandleChannel(botModel, tgMsg)
		} else {
			global.GVA_LOG.Debug("BotMsgHandlerSvc ChatGroup",
				zap.Any("msg", getUpdateText(tgMsg))) // 修复后的取文本函数
			handler.HandelChatGroup(botModel, tgMsg)
		}
	}

	return nil
}

// HandelChatGroup 处理群频道
func (handler *BotHandler) HandleChannel(botModel bot.Bot, tgMsg botapi.Update) {
	SyncChannel(botModel, tgMsg)
}

// HandelChatGroup 处理群聊消息
func (handler *BotHandler) HandelChatGroup(botModel bot.Bot, tgMsg botapi.Update) (err error) {
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

		ledger2.Handle(botModel, tgMsg)
	}()

	return
}
