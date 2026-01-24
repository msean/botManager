package listen

import (
	api "github.com/msean/botmanager/server/api"
)

type RouterGroup struct {
	ListenRoter
}

var (
	listenApi = api.ApiGroupApp.ListenApiGroup.ListenApi
)
