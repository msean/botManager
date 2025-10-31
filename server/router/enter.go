package router

import (
	"github.com/msean/botmanager/server/router/bot"
	"github.com/msean/botmanager/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System system.RouterGroup
	Bot    bot.RouterGroup
}
