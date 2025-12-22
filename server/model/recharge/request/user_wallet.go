package request

import (
	"time"

	"github.com/msean/botmanager/server/model/common/request"
)

type UserWalletSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	UserID         int64       `json:"userID" form:"userID"`
	request.PageInfo
}
