package talos

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	talosclient "github.com/siderolabs/talos/pkg/machinery/client"
)

func registerNodes(s *server.MCPServer, c *Client) {
	s.AddTool(
		mcp.NewTool("talos_get_version",
			mcp.WithDescription("Get Talos version info from one or more nodes"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.Version(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_reboot",
			mcp.WithDescription("Reboot one or more nodes"),
			mcp.WithString("nodes", mcp.Required(), mcp.Description("Comma-separated node addresses")),
			mcp.WithBoolean("powercycle", mcp.Description("Use powercycle mode instead of soft reboot (default false)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			powercycle, _ := req.Params.Arguments["powercycle"].(bool)
			ctx = c.withNodes(ctx, nodes)

			var err error
			if powercycle {
				err = c.c.Reboot(ctx, talosclient.WithPowerCycle)
			} else {
				err = c.c.Reboot(ctx)
			}
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("reboot initiated on nodes: %s", nodes)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("talos_shutdown",
			mcp.WithDescription("Shut down one or more nodes"),
			mcp.WithString("nodes", mcp.Required(), mcp.Description("Comma-separated node addresses")),
			mcp.WithBoolean("force", mcp.Description("Force shutdown even if processes are running")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			force, _ := req.Params.Arguments["force"].(bool)
			ctx = c.withNodes(ctx, nodes)

			var err error
			if force {
				err = c.c.Shutdown(ctx, talosclient.WithShutdownForce(true))
			} else {
				err = c.c.Shutdown(ctx)
			}
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("shutdown initiated on nodes: %s", nodes)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("talos_upgrade",
			mcp.WithDescription("Upgrade Talos on one or more nodes"),
			mcp.WithString("nodes", mcp.Required(), mcp.Description("Comma-separated node addresses")),
			mcp.WithString("image", mcp.Required(), mcp.Description("Installer image (e.g. ghcr.io/siderolabs/installer:v1.8.0)")),
			mcp.WithBoolean("preserve", mcp.Description("Preserve data on upgrade (default false)")),
			mcp.WithBoolean("stage", mcp.Description("Stage the upgrade and apply on next reboot")),
			mcp.WithBoolean("force", mcp.Description("Force upgrade, bypassing version checks")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			image, _ := req.Params.Arguments["image"].(string)
			preserve, _ := req.Params.Arguments["preserve"].(bool)
			stage, _ := req.Params.Arguments["stage"].(bool)
			force, _ := req.Params.Arguments["force"].(bool)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.Upgrade(ctx, image, preserve, stage, force)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)

	s.AddTool(
		mcp.NewTool("talos_get_mounts",
			mcp.WithDescription("List mounted filesystems on one or more nodes"),
			mcp.WithString("nodes", mcp.Description("Comma-separated node addresses (uses config defaults if omitted)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			nodes, _ := req.Params.Arguments["nodes"].(string)
			ctx = c.withNodes(ctx, nodes)
			resp, err := c.c.Mounts(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return protoResult(resp)
		},
	)
}
