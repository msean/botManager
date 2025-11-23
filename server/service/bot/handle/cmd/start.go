package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/service/cache"
)

func StartHandlerfunc(update tgbotapi.Update, token string, cfg cache.BotCmdCache) {
	SendCfgMessage(update, token, cfg, 2)
}
