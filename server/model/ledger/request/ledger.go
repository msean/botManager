package request

import (
	"time"

	"github.com/msean/botmanager/server/model/common/request"
)

type LedgerSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	request.PageInfo
}
