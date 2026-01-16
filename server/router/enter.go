package router

import (
	"github.com/msean/botmanager/server/router/bot"
	"github.com/msean/botmanager/server/router/public"
	"github.com/msean/botmanager/server/router/recharge"
	"github.com/msean/botmanager/server/router/system"
	"github.com/msean/botmanager/server/router/usage"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System   system.RouterGroup
	Bot      bot.RouterGroup
	Public   public.PublicRouter
	Recharge recharge.RouterGroup
	Usage    usage.RouterGroup
}
