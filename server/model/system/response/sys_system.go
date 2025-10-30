package response

import "github.com/msean/botmanager/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
