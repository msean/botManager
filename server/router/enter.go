package router

import (
	"github.com/msean/botmanager/server/router/bot"
	"github.com/msean/botmanager/server/router/example"
	"github.com/msean/botmanager/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Bot     bot.RouterGroup
}
