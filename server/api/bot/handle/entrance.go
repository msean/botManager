package handle

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/api/bot/handle/channel"
	chatgroup "github.com/msean/botmanager/server/api/bot/handle/chat_group"
	"github.com/msean/botmanager/server/api/bot/handle/private"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
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
			zap.Any("msg", getUpdateText(tgMsg)))
		private.Entrance(tgMsg, botModel)

	default:
		// 被拉入群
		if tgMsg.MyChatMember != nil {
			chatgroup.SyncChatGroup(botModel, tgMsg, nil, false)
			return nil
		}
		// 频道消息处理入口
		if tgMsg.ChannelPost != nil {
			global.GVA_LOG.Debug("BotMsgHandlerSvc  ChannelPost",
				zap.Any("msg", getUpdateText(tgMsg)))
			channel.Entrance(botModel, tgMsg)
		} else {
			// 群组消息处理入口
			global.GVA_LOG.Debug("BotMsgHandlerSvc ChatGroup",
				zap.Any("msg", getUpdateText(tgMsg)))
			chatgroup.Entrance(botModel, tgMsg)
		}
	}

	return nil
}
