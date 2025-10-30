
package request

import (
	"github.com/msean/botmanager/server/model/common/request"
	"time"
)

type BotMsgMgrSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      Ban_content  *string `json:"ban_content" form:"ban_content"` 
    request.PageInfo
}
