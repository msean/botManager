package cache

import "github.com/msean/botmanager/server/model/bot"

type BotCache struct {
	BotID       int64  `json:"botID" form:"botID"` //机器人ID
	Name        string `json:"name" form:"name"`
	IsForLedger int64  `json:"isForLedger" form:"isForLedger"` // 是否用来做记账
	Token       string `json:"token" form:"token"`             //机器人token
}

var _ CacheItem = BotCache{}

func (BotCache) TableName() string { return bot.Bot{}.TableName() }

func NewBotCache(botID int64) *BotCache {
	return &BotCache{
		BotID: botID,
	}
}

func (c BotCache) Pairs() []KvPkPair {
	return []KvPkPair{
		{"bot_id", c.BotID},
	}
}

func (BotCache) LoadType() LoadType { return LoadFromDBGet }
func (c BotCache) Release() error   { return CacheDelete(c) }
