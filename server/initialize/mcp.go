package initialize

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/msean/botmanager/server/global"
	mcpTool "github.com/msean/botmanager/server/mcp"
)

func McpRun() *server.SSEServer {
	config := global.GVA_CONFIG.MCP

	s := server.NewMCPServer(
		config.Name,
		config.Version,
	)

	global.GVA_MCP_SERVER = s

	mcpTool.RegisterAllTools(s)

	return server.NewSSEServer(s,
		server.WithSSEEndpoint(config.SSEPath),
		server.WithMessageEndpoint(config.MessagePath),
		server.WithBaseURL(config.UrlPrefix))
}
