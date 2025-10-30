package initialize

import (
	_ "github.com/msean/botmanager/server/source/example"
	_ "github.com/msean/botmanager/server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
