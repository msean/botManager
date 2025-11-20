package dao

import (
	"errors"

	"github.com/msean/botmanager/server/model/bot"
	"gorm.io/gorm"
)

type botChannelDao struct{}

func newbotChannelDao() *botChannelDao {
	return &botChannelDao{}
}

func (dao *botChannelDao) FromBotID(db *gorm.DB, channelID int) (_model bot.BotChannel, has bool, err error) {
	err = db.First(&_model, "channel_id = ?", channelID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	has = true
	return
}

func (dao *botChannelDao) MappByChannelIDList(db *gorm.DB, channelIDList []int) (mapper map[int]bot.BotChannel, err error) {
	var models []bot.BotChannel
	mapper = make(map[int]bot.BotChannel)
	if err = db.Find(&models, "channel_id in (?)", channelIDList).Error; err != nil {
		return
	}
	for _, model := range models {
		mapper[int(model.ChannelID)] = model
	}
	return
}
