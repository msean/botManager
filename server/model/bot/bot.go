// 自动生成模板Bot
package bot

import (
	"time"

	"gorm.io/gorm"
)

// 机器人 结构体  Bot
type Bot struct {
	CreatedAt time.Time      `json:"createdAt" form:"createdAt" gorm:"column:created_at;"`                        // 创建时间
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;"`                        // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                                              // 删除时间
	Name      string         `json:"name" form:"name" gorm:"comment:机器人名称;column:name;"`                          //机器人名称
	BotID     int            `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;index"`                 //机器人ID
	Token     string         `json:"token" form:"token" gorm:"comment:token;column:token;size:256;"`              //机器人token
	Chats     []BotChatGroup `json:"botChatGroups" form:"botChatGroups" gorm:"foreignKey:BotID;references:BotID"` //机器人token
	Channels  []BotChannel   `json:"botChannels" form:"botChannels" gorm:"foreignKey:BotID;references:BotID"`     //机器人token
}

// TableName 机器人 Bot自定义表名 bot
func (Bot) TableName() string {
	return "bot"
}
