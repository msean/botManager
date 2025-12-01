// 自动生成模板UserWallet
package recharge

import (
	"github.com/msean/botmanager/server/global"
)

// 用户钱包 结构体  UserWallet
type UserWallet struct {
	global.GVA_MODEL
	UserID   int64   `json:"userID" form:"userID" gorm:"comment:用户ID;column:user_id;"`      //用户ID
	UserName string  `json:"userName" form:"userName" gorm:"comment:用户名;column:user_name;"` //用户名称
	Balance  float64 `json:"balance" form:"balance" gorm:"comment:余额;column:balance;"`      //余额
	BotID    int64   `json:"botID" form:"botID" gorm:"comment:机器人名称;column:bot_id;"`        //机器人ID
}

// TableName 用户钱包 UserWallet自定义表名 bot_user_wallet
func (UserWallet) TableName() string {
	return "bot_user_wallet"
}
