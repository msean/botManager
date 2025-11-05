// 自动生成模板BotChatGroup
package bot

import (
	"github.com/msean/botmanager/server/global"
)

// 机器人群组列表 结构体  BotChatGroup
type BotChatGroup struct {
	global.GVA_MODEL
	BotName       string `json:"botName" form:"botName" gorm:"-"`                                                //机器人ID
	BotID         int64  `json:"botID" form:"botID" gorm:"column:bot_id;"`                                       //机器人ID
	ChatGroupID   int64  `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群组ID;column:chat_group_id;"`       //群组ID
	ChatGroupName string `json:"chatGroupName" form:"chatGroupName" gorm:"comment:群组ID;column:chat_group_name;"` //群组ID
}

// TableName 机器人群组列表 BotChatGroup自定义表名 bot_chat_group
func (BotChatGroup) TableName() string {
	return "bot_chat_group"
}
