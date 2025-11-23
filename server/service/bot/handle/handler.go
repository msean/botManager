package handle

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/bot/handle/cmd"
	"go.uber.org/zap"
)

type BotMsgHandlerSvc struct{}

func (svc *BotMsgHandlerSvc) Handle(c *gin.Context, botID int, body []byte) (err error) {
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

	chatType := tgMsg.Message.Chat.Type
	switch chatType {
	// 私聊
	case "private":
		global.GVA_LOG.Debug("BotMsgHandlerSvc received msg", zap.Any("msg", tgMsg.Message.Text))
		cmd.Handle(tgMsg, botModel.Token, int64(botModel.BotID))
	default:
		if tgMsg.MyChatMember != nil {
			SyncChatGroup(botModel, tgMsg)
			return nil
		}

		if tgMsg.ChannelPost != nil {
			svc.HandelChannel(botModel, tgMsg)
		} else {
			svc.HandelChatGroup(botModel, tgMsg)
		}
	}

	return nil
}

// HandelChatGroup 处理群频道
func (svc *BotMsgHandlerSvc) HandelChannel(botModel bot.Bot, tgMsg tgbotapi.Update) {
	SyncChannel(botModel, tgMsg)
}

// HandelChatGroup 处理群聊消息
func (svc *BotMsgHandlerSvc) HandelChatGroup(botModel bot.Bot, tgMsg tgbotapi.Update) (err error) {
	SyncChatGroup(botModel, tgMsg)

	// 只要消息是转发的都需要禁止
	if tgMsg.Message.ForwardFrom != nil || tgMsg.Message.ForwardFromChat != nil {
		BanUser(botModel, tgMsg, global.BanTypeForword)
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

	_, err = BanChatGroupMem(botModel, tgMsg)
	return
}

func (svc *BotMsgHandlerSvc) HandleChannel(botModel bot.Bot, tgMsg tgbotapi.Update) (err error) {
	SyncChatGroup(botModel, tgMsg)

	// 只要消息是转发的都需要禁止
	if tgMsg.Message.ForwardFrom != nil || tgMsg.Message.ForwardFromChat != nil {
		BanUser(botModel, tgMsg, global.BanTypeForword)
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

	_, err = BanChatGroupMem(botModel, tgMsg)
	return
}
