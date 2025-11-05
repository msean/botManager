package bot

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
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
	global.GVA_LOG.Info("receive telegram webhook", zap.Any("uuid", botIDStr))
	botID, err = strconv.Atoi(botIDStr)
	if err != nil {
		response.BotBadRequest(c, "invalid botID")
		return
	}

	if body, err = io.ReadAll(c.Request.Body); err != nil {
		response.BotBadRequest(c, "invalid body")
		return
	}
	c.Status(200)

	botMsgHandlerSvc.Handle(c, botID, body)
}
