
package request

import (
	"github.com/msean/botmanager/server/model/common/request"
	"time"
)

type BanRecordSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
    request.PageInfo
}
