// 自动生成模板Bot
package bot

import (
	"github.com/msean/botmanager/server/global"
)

// 机器人 结构体  Bot
type Bot struct {
	global.GVA_MODEL
	Name  *string `json:"name" form:"name" gorm:"comment:机器人名称;column:name;"`             //机器人名称
	Token *string `json:"token" form:"token" gorm:"comment:token;column:token;size:256;"` //机器人token
}

// TableName 机器人 Bot自定义表名 bot
func (Bot) TableName() string {
	return "bot"
}
