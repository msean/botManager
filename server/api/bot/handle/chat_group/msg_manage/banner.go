package msgmanage

import (
	"unicode/utf8"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func Dispatch(botModel bot.Bot, update botapi.Update, chatGroup cache.BotChatGroupCache) {
	if update.CallbackQuery != nil {
		handleCallback(botModel, &update)
		return
	} else {
		handleMsg(botModel, update, chatGroup)
	}
}

func handleMsg(botModel bot.Bot, tgMsg botapi.Update, chatGroup cache.BotChatGroupCache) (err error) {
	global.GVA_LOG.Debug("BotHandler HandelChatGroup", zap.Any("chatGroup", chatGroup))
	// 假如群聊天设置了需要禁止转发
	if chatGroup.BanForward == 1 {
		global.GVA_LOG.Debug("BotHandler HandelChatGroup", zap.Any("1", tgMsg.Message.ForwardFrom != nil), zap.Any("2", tgMsg.Message.ForwardFromChat != nil), zap.Any("2", tgMsg.Message.ExternalReply != nil))
		BanForward(botModel, tgMsg)

	}

	// 普通消息
	if tgMsg.Message == nil {
		return
	}

	// 禁用推荐联系人
	if tgMsg.Message.Contact != nil {
		return
	}

	var find bool

	// 内容检测
	if tgMsg.Message.Text != "" {
		if chatGroup.MaxWords > 0 {
			wordCount := utf8.RuneCountInString(tgMsg.Message.Text)
			if wordCount > int(chatGroup.MaxWords) {
				BanUser(botModel, tgMsg, constant.BanTypeWordLen, 0)
				return
			}
		}
		if find, err = BanChatGroupContent(botModel, tgMsg); err != nil || find {
			return
		}
	}

	global.GVA_LOG.Debug("HandelChatGroup",
		zap.Bool("find", find))

	// 成员检测
	if find, err = BanChatGroupMem(botModel, tgMsg); err != nil || find {
		return
	}

	// 检测是否关注了渠道
	ChannelFollowCheck(botModel, chatGroup.ToModel(), &tgMsg)
	return
}
