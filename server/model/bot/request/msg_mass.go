package request

import (
	"time"

	"github.com/msean/botmanager/server/model/common/request"
)

type BotMsgMassSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	request.PageInfo
}

type BotMsgMassSend struct {
	Msg string `json:"msg" binding:"required"`
	IDs []uint `json:"ids" binding:"required"`
}
