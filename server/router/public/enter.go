package public

import (
	api "github.com/msean/botmanager/server/api"
)

type RouterGroup struct {
	PublicRouter
}

var (
	medioApi = api.ApiGroupApp.PublicApiGroup.MedioApi
)
