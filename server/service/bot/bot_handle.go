package bot

import (
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

type BotHandlerSvc struct {
	botID int64
}

type AfterPublishHook func(channels []cache.BotChannelCache) error

func NewBotHandlerSvc(botID int64) *BotHandlerSvc {
	return &BotHandlerSvc{
		botID: botID,
	}
}

func (svc BotHandlerSvc) PublishAd2Channel(botApi bot_handler.Bot, chatID int64, medias []bot_handler.MediaItem, afterHook AfterPublishHook) (err error) {
	channels := cache.NewBotChannelListCache(svc.botID)
	if _, err = cache.CacheGetItem(channels); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm CacheGetItem", zap.Int64("botID", svc.botID), zap.Error(err))
		return
	}
	var hasPublishCfg bool
	cmdCfg := cache.NewBotCmdCache(svc.botID, constant.BotReplyCnfPublish2Channel, constant.BotReplyCnfType)
	if hasPublishCfg, err = cache.CacheGetItem(cmdCfg); err != nil {
		global.GVA_LOG.Error("handleBot", zap.Int("botID", int(svc.botID)), zap.Error(err))
		return
	}

	var buttons any
	if hasPublishCfg {
		buttons = bot_handler.ParseContentFromCfg(*cmdCfg, constant.ButtonTypeInline)
		global.GVA_LOG.Debug("handleBot", zap.Any("buttons", buttons))
	}
	for _, ch := range channels.Objects {
		if _, err = botApi.TgSend(ch.ChannelID, medias, buttons); err != nil {
			global.GVA_LOG.Error("HandleAdConfirm TgSend", zap.Int64("channelID", ch.ChannelID), zap.Error(err))
		}
	}
	for i := range medias {
		if medias[i].Type == "text" {
			medias[i].Text += "\n\n✅ 发布成功"
			break
		}
	}
	if _, err = botApi.TgSend(chatID, medias, nil); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", svc.botID), zap.Any("medias", medias), zap.Any("buttons", buttons), zap.Error(err))
	}

	// ✅ 执行钩子
	if afterHook != nil {
		if err = afterHook(channels.Objects); err != nil {
			global.GVA_LOG.Error("afterPublishHook 执行失败", zap.Error(err))
		}
	}
	return
}
