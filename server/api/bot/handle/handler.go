package handle

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/api/bot/handle/cmd"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"go.uber.org/zap"
)

type BotHandler struct{}

func NewBotHandler() *BotHandler {
	return &BotHandler{}
}

func getUpdateText(u tgbotapi.Update) string {
	if u.Message != nil {
		return u.Message.Text
	}
	if u.CallbackQuery != nil {
		return u.CallbackQuery.Data
	}
	return ""
}

func (handler *BotHandler) Handle(c *gin.Context, botID int, body []byte) (err error) {
	var tgMsg tgbotapi.Update
	if err = json.Unmarshal(body, &tgMsg); err != nil {
		global.GVA_LOG.Error("invalid telegram tgMsg", zap.Error(err))
		return
	}

	botModel, has, err := dao.BotDao.FromBotID(global.GVA_DB, botID)
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
	case "private":
		global.GVA_LOG.Info("BotMsgHandlerSvc received msg",
			zap.Any("msg", getUpdateText(tgMsg))) // 修复后的取文本函数

		cmd.Handle(tgMsg, botModel.Token, int64(botModel.BotID))

	default:
		if tgMsg.MyChatMember != nil {
			SyncChatGroup(botModel, tgMsg)
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
func (handler *BotHandler) HandleChannel(botModel bot.Bot, tgMsg tgbotapi.Update) {
	SyncChannel(botModel, tgMsg)
}

// HandelChatGroup 处理群聊消息
func (handler *BotHandler) HandelChatGroup(botModel bot.Bot, tgMsg tgbotapi.Update) (err error) {
	SyncChatGroup(botModel, tgMsg)

	// 只要消息是转发的都需要禁止
	if tgMsg.Message.ForwardFrom != nil || tgMsg.Message.ForwardFromChat != nil {
		BanUser(botModel, tgMsg, constant.BanTypeForword)
		return nil
	}

	// 普通消息
	if tgMsg.Message == nil {
		return nil
	}

	var find bool
	if tgMsg.Message.Text != "" {
		if find, err = BanChatGroupContent(botModel, tgMsg); err != nil || find {
			return
		}
	}

	global.GVA_LOG.Debug("HandelChatGroup",
		zap.Bool("find", find))

	_, err = BanChatGroupMem(botModel, tgMsg)
	return
}
