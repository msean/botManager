package listen

import "github.com/msean/botmanager/server/service"

type ApiGroup struct {
	ListenApi
}

var (
	listenService = service.ServiceGroupApp.ListenServiceGroup.ListenSvc
)
