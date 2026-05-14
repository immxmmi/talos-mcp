package talos

import (
	"context"
	"io"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSystem(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("talos_get_memory",
			mcp.WithDescription("Get memory usage information from one or more nodes"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.Memory(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_list_processes",
			mcp.WithDescription("List running processes on one or more nodes"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.Processes(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_get_dmesg",
			mcp.WithDescription("Get kernel ring buffer (dmesg) from one or more nodes"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			stream, err := c.c.Dmesg(ctx, false, true)
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

	s.AddTool(
		mcp.NewTool("talos_get_disks",
			mcp.WithDescription("List physical disks on one or more nodes"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.Disks(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)
}
