package system

import (
	"github.com/msean/botmanager/server/config"
)

// 配置文件结构体
type System struct {
	Config config.Server `json:"config"`
}
