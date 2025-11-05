package dao

import (
	"github.com/msean/botmanager/server/model/bot"
	"gorm.io/gorm"
)

type botChatGroupDao struct{}

func newbotChatGroupDao() *botChatGroupDao {
	return &botChatGroupDao{}
}

func (dao *botChatGroupDao) MappByChatGroupIDList(db *gorm.DB, chatGroupIDList []int) (mapper map[int]bot.BotChatGroup, err error) {
	var models []bot.BotChatGroup
	mapper = make(map[int]bot.BotChatGroup)
	if err = db.Find(&models, "chat_group_id in (?)", chatGroupIDList).Error; err != nil {
		return
	}
	for _, model := range models {
		mapper[int(model.ChatGroupID)] = model
	}
	return
}
