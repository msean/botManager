// 自动生成模板BotMsgMass
package bot

import (
	"github.com/msean/botmanager/server/global"
)

// 机器人群发 结构体  BotMsgMass
type (
	// BotMsgMass struct {
	// 	BotChatGroup
	// 	Members string `json:"members" form:"members" gorm:"comment:成员;column:members;type:text;"` //成员
	// }
	BotMassMsgRecord struct {
		global.GVA_MODEL
		BotID       int64  `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                   //机器人
		ChatGroupID int64  `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群聊ID;column:chat_group_id;"` //群聊
		Msg         string `json:"msg" form:"msg" gorm:"comment:发送消息;column:msg;type:text;"`                 //发送消息
		Members     string `json:"members" form:"members" gorm:"comment:发送成员;column:members;type:text;"`     //发送成员
		Remark      string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;type:text;"`          //发送成员
		BotFeildExtend
	}
	BotMassMsgPermission struct {
		global.GVA_MODEL
		BotName       string `json:"botName" form:"botName" gorm:"-;"`
		ChatGroupName string `json:"chatGroupName" form:"chatGroupName" gorm:"-;"`
		BotID         int64  `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id"`
		ChatGroupID   int64  `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群聊ID;column:chat_group_id;"` //群聊ID
		OprUsers      string `json:"oprUsers" form:"oprUsers" gorm:"comment:操作人;column:opr_users;type:text;"`  //操作人
	}
)

// TableName 机器人群发 BotMsgMass自定义表名 bot_msg_mass
// func (BotMsgMass) TableName() string {
// 	return "bot_chat_group"
// }

// TableName 群发历史记录 BotMassMsgRecord自定义表名 bot_mass_msg_record
func (BotMassMsgRecord) TableName() string {
	return "bot_mass_msg_record"
}

func (BotMassMsgPermission) TableName() string {
	return "bot_mass_msg_permission"
}
