// 自动生成模板BotMsgMgr
package bot

import (
	"github.com/msean/botmanager/server/global"
)

// 机器人消息管理 结构体  BotMsgMgr
type BotMsgMgr struct {
	global.GVA_MODEL
	Ban_content *string `json:"ban_content" form:"ban_content" gorm:"comment:ban_content;column:ban_content;size:1024;"` //禁用内容
	Bot_id      *int64  `json:"bot_id" form:"bot_id" gorm:"comment:机器人ID;column:bot_id;"`                                //机器人
}

// TableName 机器人消息管理 BotMsgMgr自定义表名 bot_msg_mgr
func (BotMsgMgr) TableName() string {
	return "bot_msg_mgr"
}
