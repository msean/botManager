// 自动生成模板BotMsgMass
package bot

import (
	"github.com/msean/botmanager/server/global"
)

// 机器人群发 结构体  BotMsgMass
type (
	BotMsgMass struct {
		global.GVA_MODEL
		Title       string `json:"title" form:"title" gorm:"comment:成员;column:title;size:256"`
		BotID       int64  `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                   //机器人ID
		ChatGroupID int64  `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群聊ID;column:chat_group_id;"` //群聊ID
		Members     string `json:"members" form:"members" gorm:"comment:成员;column:members;type:text;"`       //成员
		BotFeildExtend
	}
	BotMassMsgRecord struct {
		global.GVA_MODEL
		BotID       int64  `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                   //机器人
		ChatGroupID int64  `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群聊ID;column:chat_group_id;"` //群聊
		Msg         string `json:"msg" form:"msg" gorm:"comment:发送消息;column:msg;type:text;"`                 //发送消息
		Members     string `json:"members" form:"members" gorm:"comment:发送成员;column:members;type:text;"`     //发送成员
		Remark      string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;type:text;"`          //发送成员
		BotFeildExtend
	}
)

// TableName 机器人群发 BotMsgMass自定义表名 bot_msg_mass
func (BotMsgMass) TableName() string {
	return "bot_msg_mass"
}

// TableName 群发历史记录 BotMassMsgRecord自定义表名 bot_mass_msg_record
func (BotMassMsgRecord) TableName() string {
	return "bot_mass_msg_record"
}
