
package request

import (
	"github.com/msean/botmanager/server/model/common/request"
	"time"
)

type LedgerAccountGroupSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
    request.PageInfo
}
