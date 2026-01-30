package msgmanage

import (
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func BanChatGroupContent(botModel bot.Bot, tgMsg botapi.Update) (find bool, err error) {
	if tgMsg.Message == nil {
		return
	}

	cacheObjects := cache.NewBotBanContentListCache(botModel.BotID)
	_, err = cache.CacheGetItem(cacheObjects)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int64("botID", botModel.BotID), zap.Error(err))
		return
	}

	messageText := tgMsg.Message.Text

	for _, rule := range cacheObjects.Objects {
		if strings.Contains(messageText, rule.BanContent) {
			global.GVA_LOG.Info("found banned word",
				zap.String("word", rule.BanContent),
				zap.String("user", tgMsg.Message.From.UserName),
				zap.String("msg", messageText),
			)

			BanUser(botModel, tgMsg, constant.BanTypeWord, 0)

			find = true
			return
		}
	}
	return
}
