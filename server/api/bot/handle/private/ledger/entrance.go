package ledger

import (
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func Entrance(update botapi.Update, botModel bot.Bot) (err error) {

	if update.CallbackQuery != nil {
		if err = HandleCallback(update, botModel); err != nil {
			global.GVA_LOG.Error("Handle HandleCallback", zap.Any("botModel", botModel), zap.Any("callback", update), zap.Error(err))
		}
		return
	}

	// 处理正常私聊信息
	if err := Handle(botModel, update); err != nil {
		global.GVA_LOG.Error("Handle CmdParse", zap.Any("botModel", botModel), zap.Any("callback", update), zap.Error(err))
	}
	return
}

// 处理按钮回传
func HandleCallback(update botapi.Update, botModel bot.Bot) (err error) {
	return
}
