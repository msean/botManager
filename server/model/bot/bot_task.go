// 自动生成模板BotTask
package bot

import (
	"encoding/json"
	"time"

	"github.com/msean/botmanager/server/global"
)

// 任务列表 结构体  BotTask
type BotTask struct {
	global.GVA_MODEL
	Title           string          `json:"title" form:"title" gorm:"comment:发送标题;column:title;"`                                    //机器人ID
	BotID           int64           `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                                  //机器人ID
	BotName         string          `json:"botName" form:"botName" gorm:"-"`                                                         //机器人ID
	ChatGroupID     int64           `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群ID;column:chat_group_id;"`                 //群ID
	ChatGroupName   string          `json:"chatGroupName" form:"chatGroupName" gorm:"-"`                                             //机器人ID
	TaskSendType    int64           `json:"taskSendType" form:"taskSendType" gorm:"comment:发送类型;column:task_send_type;"`             //发送类型
	Content         string          `json:"content" form:"content" gorm:"column:content;type:text;"`                                 //发送内容
	ExtrendButton   json.RawMessage `json:"extrendButton" form:"extrendButton" gorm:"comment:扩展按钮;column:extrend_button;type:text;"` //扩展按钮
	SendInterval    int64           `json:"sendInterval" form:"sendInterval" gorm:"comment:发送间隔;column:send_interval;"`              //发送间隔
	NextSendTimeStr string          `json:"NextSendTimeStr" form:"NextSendTimeStr" gorm:"-"`                                         //下一次发送时间
	NextSendTime    time.Time       `json:"nextSendTime" form:"nextSendTime" gorm:"comment:下一次发送时间;column:next_send_time;"`
	StopTime        time.Time       `json:"stopTime" form:"stopTime" gorm:"comment:任务结束时间;column:stop_time;"`            //任务结束时间
	PreSendTime     *time.Time      `json:"preSendTime" form:"preSendTime" gorm:"comment:上一次发送时间;column:pre_send_time;"` //上一次发送时间
	Status          int64           `json:"status" form:"status" gorm:"comment:状态(1开 2关);column:status;"`                //状态
}

type ButtonItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// TableName 任务列表 BotTask自定义表名 bot_task
func (BotTask) TableName() string {
	return "bot_task"
}
