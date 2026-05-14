package talos

import "github.com/mark3labs/mcp-go/server"

func RegisterTools(s *server.MCPServer, c *Client) {
	registerNodes(s, c)
	registerServices(s, c)
	registerSystem(s, c)
	registerEtcd(s, c)
}
