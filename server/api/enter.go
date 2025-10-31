package v1

import (
	"github.com/msean/botmanager/server/api/bot"
	"github.com/msean/botmanager/server/api/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup system.ApiGroup
	BotApiGroup    bot.ApiGroup
}
