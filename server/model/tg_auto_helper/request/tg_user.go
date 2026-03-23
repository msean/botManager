package request

import (
	"time"

	"github.com/msean/botmanager/server/model/common/request"
)

type TgUserSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	request.PageInfo
}

type CollectGroupTaskSearch struct {
	request.PageInfo
	CreatedAtRange []string `json:"createdAtRange"`
	Status         int      `json:"status"`
	SearchText     string   `json:"searchText"`
}

// request
type CollectGroupInfoSearch struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	TaskID     uint   `json:"taskID"`
	SearchText string `json:"searchText"` // 搜索群名
}
