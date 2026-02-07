package router

import (
	"github.com/msean/botmanager/server/router/bot"
	"github.com/msean/botmanager/server/router/ledger"
	"github.com/msean/botmanager/server/router/listen"
	"github.com/msean/botmanager/server/router/public"
	"github.com/msean/botmanager/server/router/recharge"
	"github.com/msean/botmanager/server/router/system"
	"github.com/msean/botmanager/server/router/tg_auto_helper"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System         system.RouterGroup
	Bot            bot.RouterGroup
	Public         public.PublicRouter
	Recharge       recharge.RouterGroup
	Ledger         ledger.RouterGroup
	Listen         listen.RouterGroup
	Tg_auto_helper tg_auto_helper.RouterGroup
}
