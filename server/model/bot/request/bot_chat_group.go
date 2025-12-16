package request

import (
	"time"

	"github.com/msean/botmanager/server/model/common/request"
)

type BotChatGroupSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	request.PageInfo
}

// api/dto/bot_chat_message.go
type ChatMessageQuery struct {
	BotID       int64     `json:"botID" form:"botID" binding:"required"`
	ChatGroupID int64     `json:"chatGroupID" form:"chatGroupID" binding:"required"`
	UserID      int64     `json:"userId" form:"userId"`
	Username    string    `json:"username" form:"username"`
	StartTime   time.Time `json:"startTime" form:"startTime"`
	EndTime     time.Time `json:"endTime" form:"endTime"`
	Page        int       `json:"page" form:"page"`
	PageSize    int       `json:"pageSize" form:"pageSize"`
}
