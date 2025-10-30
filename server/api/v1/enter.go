package v1

import (
	"github.com/msean/botmanager/server/api/v1/bot"
	"github.com/msean/botmanager/server/api/v1/example"
	"github.com/msean/botmanager/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	BotApiGroup     bot.ApiGroup
}
