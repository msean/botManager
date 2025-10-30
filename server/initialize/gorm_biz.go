package initialize

import (
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(bot.BotMsgMgr{}, bot.Bot{})
	if err != nil {
		return err
	}
	return nil
}
