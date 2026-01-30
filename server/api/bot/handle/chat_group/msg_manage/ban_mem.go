package msgmanage

import (
	"fmt"
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

func BanChatGroupMem(botModel bot.Bot, tgMsg botapi.Update) (found bool, err error) {
	if tgMsg.Message == nil {
		return
	}

	chatID := tgMsg.Message.Chat.ID
	user := tgMsg.Message.From

	cacheObjects := cache.NewBotChatGroupBanMemCListCache(botModel.BotID, chatID)
	_, err = cache.CacheGetItem(cacheObjects)
	global.GVA_LOG.Info("BanChatGroupMem",
		zap.Any("cacheObjects", cacheObjects),
		zap.Error(err),
	)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int64("botID", botModel.BotID), zap.Error(err))
		return
	}
	fullName := fmt.Sprintf("%s%s", user.FirstName, user.LastName)

	for _, ban := range cacheObjects.Objects {
		banStr := strings.ToLower(strings.TrimSpace(ban.BanMemContent))
		if banStr == "" {
			continue
		}

		if strings.Contains(fullName, banStr) {
			global.GVA_LOG.Info("found banned member",
				zap.String("ban_word", ban.BanMemContent),
				zap.String("member", user.UserName),
				zap.Int64("chatID", chatID),
			)

			BanUser(botModel, tgMsg, constant.BanTypeMem, 0)

			found = true
			return
		}
	}
	return
}
