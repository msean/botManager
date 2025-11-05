package dao

import (
	"github.com/msean/botmanager/server/model/bot"
	"gorm.io/gorm"
)

type botGroupMemDao struct{}

func newbotGroupMemDao() *botGroupMemDao {
	return &botGroupMemDao{}
}

func (dao *botGroupMemDao) ListByBotIDAndChatGroupID(db *gorm.DB, botID int, chatGroupID int) (members []bot.BotBanGroupMem, err error) {
	err = db.Find(&members, "bot_id = ? and chat_group_id=?", botID, chatGroupID).Error
	return
}
