package ast

import (
	"path/filepath"

	"github.com/msean/botmanager/server/global"
)

func init() {
	global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("../../../")
	global.GVA_CONFIG.AutoCode.Server = "server"
}
