package talos

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	machineapi "github.com/siderolabs/talos/pkg/machinery/api/machine"
)

func registerEtcd(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("talos_etcd_members",
			mcp.WithDescription("List etcd cluster members"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.EtcdMemberList(ctx, &machineapi.EtcdMemberListRequest{})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_etcd_status",
			mcp.WithDescription("Get etcd member status (alarm, size, hash, leader info)"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.EtcdStatus(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_etcd_alarm_list",
			mcp.WithDescription("List active etcd alarms"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.EtcdAlarmList(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_etcd_defragment",
			mcp.WithDescription("Defragment etcd storage on one or more nodes"),
			mcp.WithString("nodes", mcp.Required(), mcp.Description("Comma-separated node addresses")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.EtcdDefragment(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)
}
