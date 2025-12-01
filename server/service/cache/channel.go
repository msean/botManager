package cache

import (
	"fmt"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/utils"
	"go.uber.org/zap"
)

type (
	BotChannelCache struct {
		BotID       int64  `json:"botID"`
		ChannelID   int64  `json:"channelID"`
		ChannelName string `json:"channelName"`
	}
	BotChannelListCache struct {
		BotID   int64             `json:"botID"`
		Objects []BotChannelCache `json:"objects"`
	}
)

func (c BotChannelCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
		{"channel_id", c.ChannelID},
	}
}
func (c BotChannelListCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
	}
}

func (BotChannelCache) LoadType() LoadType     { return LoadFromDBGet }
func (BotChannelListCache) LoadType() LoadType { return LoadFromDBList }

func (BotChannelCache) TableName() string     { return bot.BotChannel{}.TableName() }
func (BotChannelListCache) TableName() string { return bot.BotChannel{}.TableName() }

func NewBotChannelCache(botID int64, channelID int64) *BotChannelCache {
	return &BotChannelCache{
		BotID:     botID,
		ChannelID: channelID,
	}
}

func NewBotChannelListCache(botID int64) *BotChannelListCache {
	return &BotChannelListCache{
		BotID: botID,
	}
}

func (c BotChannelListCache) Release() error { return CacheDelete(c) }
func (c BotChannelCache) Release() error     { return CacheDelete(c) }

func ReleaseChannelModelChange(modelID uint) (err error) {
	var object bot.BotChannel
	var has bool
	if has, err = utils.Get(global.GVA_DB, &object, utils.IDCond(modelID)); !has || err != nil {
		if !has {
			err = fmt.Errorf("record not found")
		}
		global.GVA_LOG.Error("ReleaseRechargeCnf", zap.Any("id", modelID), zap.Error(err))
		return
	}

	if err = NewBotChannelListCache(object.BotID).Release(); err != nil {
		global.GVA_LOG.Error("ReleaseRechargeCnf", zap.Any("object", object), zap.Error(err))
	}
	if err = NewBotChannelCache(object.BotID, object.ChannelID).Release(); err != nil {
		global.GVA_LOG.Error("ReleaseRechargeCnf", zap.Any("object", object), zap.Error(err))
	}
	return
}
