package bot

import (
	"context"

	"github.com/msean/botmanager/server/model/bot"
)

type BotMsgHandlerSvc struct{}

func (BotMsgHandlerSvc *BotMsgHandlerSvc) CreateBotMsgMgr(ctx context.Context, bot_msg_mgr *bot.BotBanContent) (err error) {
	// var bots
	// err = global.GVA_DB.(bot_msg_mgr).Error
	return err
}
