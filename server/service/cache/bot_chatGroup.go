package cache

import (
	"fmt"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/utils"
	"go.uber.org/zap"
)

type (
	BotChatGroupCache struct {
		BotID         int64  `json:"botID"`
		ChatGroupID   int64  `json:"chatGroupID"`
		SyncMessage   int64  `json:"syncMessage"` // 是否需要同步消息
		BanForward    int64  `json:"banForward"`  // 是否禁用转发消息
		ChatGroupName string `json:"chatGroupName"`
	}
)

func NewBotChatGroupCache(botID, chatGroupID int64) *BotChatGroupCache {
	return &BotChatGroupCache{
		BotID:       botID,
		ChatGroupID: chatGroupID,
	}
}

func (BotChatGroupCache) TableName() string { return bot.BotChatGroup{}.TableName() }

func (c BotChatGroupCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
		{"chat_group_id", c.ChatGroupID},
	}
}

func (BotChatGroupCache) LoadType() LoadType { return LoadFromDBGet }
func (c BotChatGroupCache) Release() error   { return CacheDelete(c) }

func ReleaseBotChatGroup(modelID int) (err error) {
	var object bot.BotChatGroup
	var has bool
	if has, err = utils.Get(global.GVA_MYSQL, &object, utils.IDCond(modelID)); !has || err != nil {
		if !has {
			err = fmt.Errorf("record not found")
		}
		global.GVA_LOG.Error("ReleaseRechargeCnf", zap.Any("id", modelID), zap.Error(err))
		return
	}

	if err = NewBotChatGroupCache(object.BotID, object.ChatGroupID).Release(); err != nil {
		global.GVA_LOG.Error("ReleaseRechargeCnf", zap.Any("id", modelID), zap.Error(err))
	}
	return
}
