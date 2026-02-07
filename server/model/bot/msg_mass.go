// 自动生成模板BotMsgMass
package bot

import (
	"github.com/msean/botmanager/server/global"
)

// 机器人群发 结构体  BotMsgMass
type BotMsgMass struct {
	global.GVA_MODEL
	BotID       int64  `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                   //机器人ID
	ChatGroupID int64  `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群聊ID;column:chat_group_id;"` //群聊ID
	Members     string `json:"members" form:"members" gorm:"comment:成员;column:members;type:text;"`       //成员
	BotFeildExtend
}

// TableName 机器人群发 BotMsgMass自定义表名 bot_msg_mass
func (BotMsgMass) TableName() string {
	return "bot_msg_mass"
}
