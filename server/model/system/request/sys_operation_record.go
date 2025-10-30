package request

import (
	"github.com/msean/botmanager/server/model/common/request"
	"github.com/msean/botmanager/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
