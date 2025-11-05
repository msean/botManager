package cache

import (
	"github.com/msean/botmanager/server/model/bot"
)

type (
	BotChatGroupBanMemCache struct {
		BanMemContent string `json:"banMemContent"`
	}
	BotChatGroupCache struct {
		ChatGroupID   int64  `json:"chatGroupID" form:"chatGroupID"`     //群组ID
		ChatGroupName string `json:"chatGroupName" form:"chatGroupName"` //群组ID
		BotID         int    `json:"botID" form:"botID"`                 //机器人ID
	}
	BotBanContentCache struct {
		BanContent string `json:"banContent" form:"banContent" ` //禁用内容
		BotID      int64  `json:"botID" form:"botID" `           //机器人ID
	}
)

func (cache BotChatGroupBanMemCache) TableName() string {
	return bot.BotBanGroupMem{}.TableName()
}

func (cache BotChatGroupCache) TableName() string {
	return bot.BotChatGroup{}.TableName()
}

func (cache BotBanContentCache) TableName() string {
	return bot.BotBanContent{}.TableName()
}

var _ bot.BotChatGroup

func BotChatGroupMemPk(botID, chatGroupID int) []KvPkPair {

	return []KvPkPair{
		{
			PKCol: "bot_id",
			PKVal: botID,
		},
		{
			PKCol: "chat_group_id",
			PKVal: chatGroupID,
		},
	}
}

func BotChatGroupPk(botID, chatGroupID int) []KvPkPair {
	return []KvPkPair{
		{
			PKCol: "bot_id",
			PKVal: botID,
		},
		{
			PKCol: "chat_group_id",
			PKVal: chatGroupID,
		},
	}
}

func BotBanContentPk(botID int) []KvPkPair {
	return []KvPkPair{
		{
			PKCol: "bot_id",
			PKVal: botID,
		},
	}
}

func ReleaseBotBanContent(botID int) (err error) {
	return CacheDelete(BotBanContentCache{}.TableName(), BotBanContentPk(botID))
}

func ReleaseBotChatGroup(botID, chatGroupID int) (err error) {
	return CacheDelete(BotChatGroupCache{}.TableName(), BotChatGroupPk(botID, chatGroupID))
}

func ReleaseBotChatGroupMem(botID, chatGroupID int) (err error) {
	return CacheDelete(BotChatGroupBanMemCache{}.TableName(), BotChatGroupPk(botID, chatGroupID))
}
