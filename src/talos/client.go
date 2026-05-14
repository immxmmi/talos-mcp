package talos

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	talosclient "github.com/siderolabs/talos/pkg/machinery/client"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/immxmmi/talos-mcp/config"
)

type Client struct {
	c            *talosclient.Client
	defaultNodes []string
	perms        config.Permissions
}

func NewClient(cfg *config.Config) (*Client, error) {
	var opts []talosclient.OptionFunc

	if cfg.TalosConfig != "" {
		opts = append(opts, talosclient.WithConfigFromFile(cfg.TalosConfig))
	} else {
		opts = append(opts, talosclient.WithDefaultConfig())
	}

	if cfg.Context != "" {
		opts = append(opts, talosclient.WithContextName(cfg.Context))
	}

	if len(cfg.Endpoints) > 0 {
		opts = append(opts, talosclient.WithEndpoints(cfg.Endpoints...))
	}

	c, err := talosclient.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		c:            c,
		defaultNodes: cfg.Nodes,
		perms:        cfg.Perms,
	}, nil
}

func (c *Client) guardExec() (*mcp.CallToolResult, bool) {
	if !c.perms.AllowExec {
		return mcp.NewToolResultError("exec/action operations are disabled (set MCP_ALLOW_EXEC=true or unset MCP_READ_ONLY)"), false
	}
	return nil, true
}

func (c *Client) Close() {
	_ = c.c.Close()
}

// withNodes returns a context targeting the given comma-separated nodes, or the configured defaults.
func (c *Client) withNodes(ctx context.Context, nodes string) context.Context {
	var targets []string
	if nodes != "" {
		for _, n := range strings.Split(nodes, ",") {
			if t := strings.TrimSpace(n); t != "" {
				targets = append(targets, t)
			}
		}
	} else {
		targets = c.defaultNodes
	}
	if len(targets) > 0 {
		return talosclient.WithNodes(ctx, targets...)
	}
	return ctx
}

// Ping validates the connection by fetching the Talos version from the configured nodes.
func (c *Client) Ping() error {
	ctx := c.withNodes(context.Background(), "")
	if _, err := c.c.Version(ctx); err != nil {
		return fmt.Errorf("Talos connection test failed: %w", err)
	}
	return nil
}

func protoResult(msg proto.Message) (*mcp.CallToolResult, error) {
	b, err := protojson.MarshalOptions{EmitUnpopulated: false}.Marshal(msg)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}
