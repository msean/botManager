package cache

import (
	"fmt"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/utils"
	"go.uber.org/zap"
)

type (
	BotChatGroupBanMemCache struct {
		BanMemContent string `json:"banMemContent" gorm:"column:banMemContent"`
	}
	BotChatGroupBanMemCListCache struct {
		BotID       int64                     `json:"botID"`
		ChatGroupID int64                     `json:"chatGroupID"`
		Objects     []BotChatGroupBanMemCache `json:"objects"`
	}
	BotChatGroupCache struct {
		BotID         int64  `json:"botID"`
		ChatGroupID   int64  `json:"chatGroupID"`
		ChatGroupName string `json:"chatGroupName"`
	}
	BotBanContentCache struct {
		BotID      int64  `json:"botID"`
		BanContent string `json:"banContent"`
	}
	BotBanContentListCache struct {
		BotID   int64                `json:"botID"`
		Objects []BotBanContentCache `json:"objects"`
	}
	BotChannelCache struct {
		BotID       int64  `json:"botID"`
		ChannelID   int64  `json:"channelID"`
		ChannelName string `json:"channelName"`
	}
	BotCmdCache struct {
		BotID      int64  `json:"botID"`
		Cmd        string `json:"cmd"`
		Type       int    `json:"type"`
		Content    string `json:"content"`
		CmdButtons string `json:"cmdButtons"`
	}
	BotCmdCacheWithNoContent struct {
		Cmd        string `json:"cmd"`
		CmdButtons string `json:"cmdButtons"`
	}
	// 获取
	BotCmdCacheList struct {
		BotID   int64                      `json:"botID"`
		Objects []BotCmdCacheWithNoContent `json:"objects"`
	}
	RechargeCnfObj struct {
		PublishTimes int     `json:"publishTimes"` //发布次数
		Price        float64 `json:"price"`        //价格
	}
	// 获取
	RechargeCnfCacheList struct {
		BotID   int64            `json:"botID"`
		Objects []RechargeCnfObj `json:"objects"`
	}
)

func NewBotChatGroupCache(object bot.BotChatGroup) *BotChatGroupCache {
	return &BotChatGroupCache{
		BotID:       object.BotID,
		ChatGroupID: object.ChatGroupID,
	}
}

func NewBotBanContentListCache(botID int64) *BotBanContentListCache {
	return &BotBanContentListCache{
		BotID: botID,
	}
}

func NewBotChannelCache(object bot.BotChannel) *BotChannelCache {
	return &BotChannelCache{
		BotID:     object.BotID,
		ChannelID: object.ChannelID,
	}
}

func NewBotCmdCache(botID int64, cmd string, _type int) *BotCmdCache {
	return &BotCmdCache{
		BotID: botID,
		Cmd:   cmd,
		Type:  _type,
	}
}

func NewBotCmdCacheList(botID int64) *BotCmdCacheList {
	return &BotCmdCacheList{
		BotID: botID,
	}
}

func NewBotChatGroupBanMemCListCache(botID int64, chatGroupID int64) *BotChatGroupBanMemCListCache {
	return &BotChatGroupBanMemCListCache{
		BotID:       botID,
		ChatGroupID: chatGroupID,
	}
}

func NewRechargeCnfListCache(botID int64) *RechargeCnfCacheList {
	return &RechargeCnfCacheList{
		BotID: botID,
	}
}

func (BotChatGroupBanMemCListCache) TableName() string { return bot.BotBanGroupMem{}.TableName() }
func (BotChatGroupCache) TableName() string            { return bot.BotChatGroup{}.TableName() }
func (BotBanContentListCache) TableName() string       { return bot.BotBanContent{}.TableName() }
func (BotChannelCache) TableName() string              { return bot.BotChannel{}.TableName() }
func (BotCmdCache) TableName() string                  { return bot.BotCmdConfig{}.TableName() }
func (BotCmdCacheList) TableName() string              { return bot.BotCmdConfig{}.TableName() }
func (RechargeCnfCacheList) TableName() string         { return recharge.RechargeConfig{}.TableName() }

func (c BotChatGroupBanMemCListCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
		{"chat_group_id", c.ChatGroupID},
	}
}

func (c BotChatGroupCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
		{"chat_group_id", c.ChatGroupID},
	}
}

func (c BotBanContentListCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
	}
}

func (c BotChannelCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
		{"channel_id", c.ChannelID},
	}
}

func (c BotCmdCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
		{"cmd", c.Cmd},
		{"type", c.Type},
	}
}

func (c BotCmdCacheList) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
		{"type", global.BotReplyCmdType},
	}
}

func (c RechargeCnfCacheList) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
	}
}

func (BotChatGroupBanMemCListCache) LoadType() LoadType { return LoadFromDBList }
func (BotChatGroupCache) LoadType() LoadType            { return LoadFromDBGet }
func (BotBanContentListCache) LoadType() LoadType       { return LoadFromDBList }
func (BotChannelCache) LoadType() LoadType              { return LoadFromDBGet }
func (BotCmdCache) LoadType() LoadType                  { return LoadFromDBGet }
func (BotCmdCacheList) LoadType() LoadType              { return LoadFromDBList }
func (RechargeCnfCacheList) LoadType() LoadType         { return LoadFromDBList }
func (SysCnfCache) LoadType() LoadType                  { return LoadFromDBGet }

func (c BotChatGroupBanMemCListCache) Release() error { return CacheDelete(c) }
func (c BotChatGroupCache) Release() error            { return CacheDelete(c) }
func (c BotBanContentListCache) Release() error       { return CacheDelete(c) }
func (c BotChannelCache) Release() error              { return CacheDelete(c) }
func (c BotCmdCache) Release() error                  { return CacheDelete(c) }
func (c BotCmdCacheList) Release() error              { return CacheDelete(c) }
func (c RechargeCnfCacheList) Release() error         { return CacheDelete(c) }

func ReleaseRechargeCnf(modelID int) (err error) {
	var object recharge.RechargeConfig
	var has bool
	if has, err = utils.Get(global.GVA_DB, &object, utils.IDCond(modelID)); !has || err != nil {
		if !has {
			err = fmt.Errorf("record not found")
		}
		global.GVA_LOG.Error("ReleaseRechargeCnf", zap.Any("id", modelID), zap.Error(err))
		return
	}

	if err = NewRechargeCnfListCache(object.BotID).Release(); err != nil {
		global.GVA_LOG.Error("ReleaseRechargeCnf", zap.Any("id", modelID), zap.Error(err))
	}
	return
}
