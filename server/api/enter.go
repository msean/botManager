package v1

import (
	"github.com/msean/botmanager/server/api/bot"
	"github.com/msean/botmanager/server/api/ledger"
	"github.com/msean/botmanager/server/api/listen"
	"github.com/msean/botmanager/server/api/public"
	"github.com/msean/botmanager/server/api/recharge"
	"github.com/msean/botmanager/server/api/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	BotApiGroup      bot.ApiGroup
	PublicApiGroup   public.ApiGroup
	RechargeApiGroup recharge.ApiGroup
	LedgerApiGroup   ledger.ApiGroup
	ListenApiGroup   listen.ApiGroup
}
