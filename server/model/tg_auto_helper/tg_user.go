// 自动生成模板TgUser
package tg_auto_helper

import (
	"github.com/msean/botmanager/server/global"
)

// telegram用户管理 结构体  TgUser
type TgUser struct {
	global.GVA_MODEL
	NickName         string `json:"nickName" form:"nickName" gorm:"comment:昵称;column:nick_name;size:128;"` //昵称
	Phone            string `json:"phone" form:"phone" gorm:"comment:手机号码;column:phone;size:64;"`          //手机号码
	PhoneCodeHash    string `json:"-" gorm:"column:phone_code_hash;size:128"`
	ApiId            int    `json:"apiId" form:"apiId" gorm:"comment:apiId;column:apiId;size:256;"`                                  //apiId
	ApiHash          string `json:"apiHash" form:"apiHash" gorm:"comment:apiHash;column:apiHash;size:512;"`                          //apiHash
	NextVerification string `json:"nextVerification" form:"nextVerification" gorm:"comment:二步验证;column:next_verification;size:216;"` //密码验证
	Status           int64  `json:"status" form:"status" gorm:"comment:状态;column:status;"`                                           //状态
	CodeVerfication  string `json:"codeVerfication" form:"codeVerfication" gorm:"-"`
	SessionPath      string `gorm:"column:session_path;size:256"`
}

// TableName telegram用户管理 TgUser自定义表名 tg_user
func (TgUser) TableName() string {
	return "tg_user"
}
