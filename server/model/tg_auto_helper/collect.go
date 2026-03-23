package tg_auto_helper

import "github.com/msean/botmanager/server/global"

type CollectGroupTask struct {
	global.GVA_MODEL
	TgUserID      int    `json:"tgUserID" gorm:"column:tg_user_id;comment:采集使用的TG账号"`
	SearchText    string `json:"searchText" gorm:"column:search_text;size:128;comment:关键词"`
	CurrentPage   int    `json:"currentPage" gorm:"column:current_page;default:1;comment:当前采集页"`
	TotalPage     int    `json:"totalPage" gorm:"column:total_page;comment:总页数"`
	TotalCount    int    `json:"totalCount" gorm:"column:total_count;comment:总群数"`
	SourceBotName string `json:"sourceBotName" gorm:"column:sourceBotName;default:soso;comment:来源"`
	Status        int    `json:"status" gorm:"column:status;default:1;comment:1采集中 2完成 3停止"`
	MemberLower   int    `json:"members_lower" gorm:"column:members_lower;default:0;comment:群人数下限 0 表示无限制"`
	Remark        string `json:"remark" gorm:"column:remark;comment:备注"`
}

type CollectGroupInfo struct {
	global.GVA_MODEL
	Title     string `json:"title" gorm:"column:title;size:128;comment:群名称"`
	GroupName string `json:"groupName" gorm:"column:group_name;size:128;index;comment:群名"`
	Link      string `json:"link" gorm:"column:link;size:256;comment:群链接"`
	Members   int    `json:"members" gorm:"column:members;comment:成员数"`
	TaskID    uint   `json:"taskID" gorm:"column:task_id;index;comment:任务ID"`
}

type CollectUser struct {
	global.GVA_MODEL
	GroupID  int64  `json:"groupID" gorm:"column:group_id;size:128;comment:groupID"`
	UserID   int64  `json:"userID" gorm:"column:user_id;size:128;comment:userID"`
	Username string `json:"username" gorm:"column:username;size:128;comment:userName"`
	NickName string `json:"nickname" gorm:"column:nickname;size:128;comment:nickName"`
}

// TableName telegram用户管理 TgUser自定义表名 tg_user
func (CollectGroupTask) TableName() string {
	return "tg_collect_group_task"
}

// TableName telegram用户管理 TgUser自定义表名 tg_user
func (CollectGroupInfo) TableName() string {
	return "tg_collect_group_info"
}

// TableName CollectUser TgUser自定义表名 tg_user
func (CollectUser) TableName() string {
	return "tg_collect_user"
}
