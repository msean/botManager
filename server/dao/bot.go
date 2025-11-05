package dao

import (
	"errors"

	"github.com/msean/botmanager/server/model/bot"
	"gorm.io/gorm"
)

type botDao struct{}

func newBotDao() *botDao {
	return &botDao{}
}

func (dao *botDao) FromBotID(db *gorm.DB, botID int) (botModel bot.Bot, has bool, err error) {
	err = db.First(&botModel, "bot_id = ?", botID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	has = true
	return
}

func (dao *botDao) ListBotBannerContentByID(db *gorm.DB, botID int) (banContents []bot.BotBanContent, err error) {
	err = db.Find(&banContents, "bot_id = ?", botID).Error
	return
}

func (dao *botDao) MappByIDList(db *gorm.DB, botIDList []int) (mapper map[int]bot.Bot, err error) {
	var bots []bot.Bot
	mapper = make(map[int]bot.Bot)
	if err = db.Find(&bots, "bot_id in (?)", botIDList).Error; err != nil {
		return
	}
	for _, botModel := range bots {
		mapper[botModel.BotID] = botModel
	}
	return
}

func (dao *botDao) All(db *gorm.DB) (bots []bot.Bot, err error) {
	err = db.Find(&bots).Error
	return
}

func (dao *botDao) AllWithChatGroup(db *gorm.DB) (bots []bot.Bot, err error) {
	err = db.Model(&bot.Bot{}).Preload("Chats").Find(&bots).Error
	return
}
