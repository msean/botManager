package cache

import (
	"github.com/msean/botmanager/server/model/bot"
)

type (
	BotChatGroupBanMemCache struct {
		BotID         int    `json:"botID"`
		ChatGroupID   int    `json:"chatGroupID"`
		BanMemContent string `json:"banMemContent"`
	}

	BotChatGroupBanMemCListCache struct {
		BotID   int                       `json:"botID"`
		Objects []BotChatGroupBanMemCache `json:"objects"`
	}

	BotChatGroupCache struct {
		BotID         int    `json:"botID"`
		ChatGroupID   int    `json:"chatGroupID"`
		ChatGroupName string `json:"chatGroupName"`
	}

	BotBanContentCache struct {
		BotID      int    `json:"botID"`
		BanContent string `json:"banContent"`
	}

	BotBanContentListCache struct {
		BotID   int                  `json:"botID"`
		Objects []BotBanContentCache `json:"objects"`
	}

	BotChannelCache struct {
		BotID       int    `json:"botID"`
		ChannelID   int    `json:"channelID"`
		ChannelName string `json:"channelName"`
	}
	BotCmdCache struct {
		BotID      int64  `json:"botID"`
		Cmd        string `json:"cmd"`
		Content    string `json:"content"`
		CmdButtons string `json:"cmdButtons"`
	}
	BotCmdCacheWithNoContent struct {
		Cmd        string `json:"cmd"`
		CmdButtons string `json:"cmdButtons"`
	}
	BotCmdCacheList struct {
		BotID   int                        `json:"botID"`
		Objects []BotCmdCacheWithNoContent `json:"objects"`
	}
)

func NewBotChatGroupBanMemCache(object bot.BotBanGroupMem) *BotChatGroupBanMemCache {
	return &BotChatGroupBanMemCache{
		BotID:       int(object.BotID),
		ChatGroupID: int(object.ChatGroupID),
	}
}

func NewBotChatGroupCache(object bot.BotChatGroup) *BotChatGroupCache {
	return &BotChatGroupCache{
		BotID:       int(object.BotID),
		ChatGroupID: int(object.ChatGroupID),
	}
}

func NewBotBanContentListCache(botID int) *BotBanContentListCache {
	return &BotBanContentListCache{
		BotID: botID,
	}
}

func NewBotChannelCache(object bot.BotChannel) *BotChannelCache {
	return &BotChannelCache{
		BotID:     int(object.BotID),
		ChannelID: int(object.ChannelID),
	}
}

func NewBotCmdCache(botID int64, cmd string) *BotCmdCache {
	return &BotCmdCache{
		BotID: botID,
		Cmd:   cmd,
	}
}

func NewBotCmdCacheList(botID int) *BotCmdCacheList {
	return &BotCmdCacheList{
		BotID: botID,
	}
}

func NewBotChatGroupBanMemCListCache(botID int) *BotChatGroupBanMemCListCache {
	return &BotChatGroupBanMemCListCache{
		BotID: botID,
	}
}

func (BotChatGroupBanMemCListCache) TableName() string { return bot.BotChatGroup{}.TableName() }
func (BotChatGroupCache) TableName() string            { return bot.BotChatGroup{}.TableName() }
func (BotBanContentListCache) TableName() string       { return bot.BotBanContent{}.TableName() }
func (BotChannelCache) TableName() string              { return bot.BotChannel{}.TableName() }
func (BotCmdCache) TableName() string                  { return bot.BotCmdConfig{}.TableName() }
func (BotCmdCacheList) TableName() string              { return bot.BotCmdConfig{}.TableName() }

func (c BotChatGroupBanMemCListCache) Pairs() []KvPkPair {
	return []KvPkPair{{"bot_id", c.BotID}}
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
	}
}

func (c BotCmdCacheList) Pairs() []KvPkPair {
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

func (c BotChatGroupBanMemCListCache) Release() error { return CacheDelete(c) }
func (c BotChatGroupCache) Release() error            { return CacheDelete(c) }
func (c BotBanContentListCache) Release() error       { return CacheDelete(c) }
func (c BotChannelCache) Release() error              { return CacheDelete(c) }
func (c BotCmdCache) Release() error                  { return CacheDelete(c) }
func (c BotCmdCacheList) Release() error              { return CacheDelete(c) }
