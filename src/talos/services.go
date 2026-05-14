package talos

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	commonapi "github.com/siderolabs/talos/pkg/machinery/api/common"
)

func registerServices(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("talos_list_services",
			mcp.WithDescription("List all services on one or more nodes"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.ServiceList(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_start_service",
			mcp.WithDescription("Start a service on one or more nodes"),
			mcp.WithString("nodes", mcp.Required(), mcp.Description("Comma-separated node addresses")),
			mcp.WithString("service", mcp.Required(), mcp.Description("Service name (e.g. kubelet, etcd)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if r, ok := c.guardExec(); !ok { return r, nil }
			nodes, _ := req.Params.Arguments["nodes"].(string)
			service, _ := req.Params.Arguments["service"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.ServiceStart(ctx, service)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_stop_service",
			mcp.WithDescription("Stop a service on one or more nodes"),
			mcp.WithString("nodes", mcp.Required(), mcp.Description("Comma-separated node addresses")),
			mcp.WithString("service", mcp.Required(), mcp.Description("Service name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if r, ok := c.guardExec(); !ok { return r, nil }
			nodes, _ := req.Params.Arguments["nodes"].(string)
			service, _ := req.Params.Arguments["service"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.ServiceStop(ctx, service)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_restart_service",
			mcp.WithDescription("Restart a service on one or more nodes"),
			mcp.WithString("nodes", mcp.Required(), mcp.Description("Comma-separated node addresses")),
			mcp.WithString("service", mcp.Required(), mcp.Description("Service name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if r, ok := c.guardExec(); !ok { return r, nil }
			nodes, _ := req.Params.Arguments["nodes"].(string)
			service, _ := req.Params.Arguments["service"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.ServiceRestart(ctx, service)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_get_service_logs",
			mcp.WithDescription("Get logs for a system service on a node"),
			mcp.WithString("nodes", mcp.Required(), mcp.Description("Node address (single node)")),
			mcp.WithString("service", mcp.Required(), mcp.Description("Service name (e.g. kubelet, etcd, containerd)")),
			mcp.WithString("tail_lines", mcp.Description("Number of log lines to return (default 100)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			service, _ := req.Params.Arguments["service"].(string)
			tailStr, _ := req.Params.Arguments["tail_lines"].(string)
			tailLines := int32(100)
			if tailStr != "" {
				var n int32
				fmt.Sscanf(tailStr, "%d", &n)
				if n > 0 {
					tailLines = n
				}
			}
			ctx = c.withNodes(ctx, nodes)
			stream, err := c.c.Logs(ctx, "system", commonapi.ContainerDriver_CONTAINERD, service, false, tailLines)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			var buf strings.Builder
			for {
				msg, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					break
				}
				buf.Write(msg.GetBytes())
			}
			return mcp.NewToolResultText(buf.String()), nil
		},
	)
}
