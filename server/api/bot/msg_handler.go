package bot

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/api/bot/handle"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/common/response"
	"go.uber.org/zap"
)

type BotMsgHandler struct{}

func (api *BotMsgHandler) Handle(c *gin.Context) {
	var body []byte
	var botID int
	var err error

	botIDStr := c.Param("botUUID")
	botID, err = strconv.Atoi(botIDStr)
	if err != nil {
		global.GVA_LOG.Error("receive telegram webhook", zap.Any("botIDStr", botIDStr))
		response.BotBadRequest(c, "invalid botID")
		return
	}

	if body, err = io.ReadAll(c.Request.Body); err != nil {
		response.BotBadRequest(c, "invalid body")
		return
	}
	c.Status(200)
	c.Writer.WriteHeaderNow()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				global.GVA_LOG.Error("telegram webhook panic", zap.Any("recover", r))
			}
		}()
		handle.NewBotHandler().Handle(c, botID, body)
	}()
}
