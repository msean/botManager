package dao

import (
	"errors"

	"github.com/msean/botmanager/server/model/bot"
	"gorm.io/gorm"
)

type botChatGroupDao struct{}

func newbotChatGroupDao() *botChatGroupDao {
	return &botChatGroupDao{}
}

func (dao *botChatGroupDao) MappByChatGroupIDList(db *gorm.DB, chatGroupIDList []int64) (mapper map[int64]bot.BotChatGroup, err error) {
	var models []bot.BotChatGroup
	mapper = make(map[int64]bot.BotChatGroup)
	if err = db.Find(&models, "chat_group_id in (?)", chatGroupIDList).Error; err != nil {
		return
	}
	for _, model := range models {
		mapper[model.ChatGroupID] = model
	}
	return
}

func (dao *botChatGroupDao) FromID(db *gorm.DB, botChatGroupID int) (botChatGroupModel bot.BotChatGroup, has bool, err error) {
	err = db.First(&botChatGroupModel, "chat_group_id = ?", botChatGroupID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	has = true
	return
}
