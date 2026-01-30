package msgmanage

import (
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func BanForward(botModel bot.Bot, tgMsg botapi.Update) {
	if tgMsg.Message.ForwardFrom != nil || tgMsg.Message.ForwardFromChat != nil || tgMsg.Message.ExternalReply != nil {
		BanUser(botModel, tgMsg, constant.BanTypeForword, 0)
		return
	}
}
