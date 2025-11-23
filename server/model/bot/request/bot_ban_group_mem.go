package request

import (
	"time"

	"github.com/msean/botmanager/server/model/common/request"
)

type BotBanGroupMemSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	BotID          int         `json:"botID" form:"botID"`
	request.PageInfo
}
