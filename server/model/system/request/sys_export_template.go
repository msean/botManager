package request

import (
	"time"

	"github.com/msean/botmanager/server/model/common/request"
	"github.com/msean/botmanager/server/model/system"
)

type SysExportTemplateSearch struct {
	system.SysExportTemplate
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}
