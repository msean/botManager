
package request

import (
	"github.com/msean/botmanager/server/model/common/request"
	"time"
)

type BotChatGroupSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
    request.PageInfo
}
