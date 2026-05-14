package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	TalosConfig string
	Context     string
	Endpoints   []string
	Nodes       []string
	Perms       Permissions
}

func Load() (*Config, error) {
	cfg := &Config{
		TalosConfig: os.Getenv("TALOSCONFIG"),
		Context:     os.Getenv("TALOS_CONTEXT"),
	}

	if cfg.TalosConfig == "" {
		home, _ := os.UserHomeDir()
		cfg.TalosConfig = filepath.Join(home, ".talos", "config")
	}

	if v := os.Getenv("TALOS_ENDPOINTS"); v != "" {
		for _, ep := range strings.Split(v, ",") {
			ep = strings.TrimSpace(ep)
			if ep != "" {
				cfg.Endpoints = append(cfg.Endpoints, ep)
			}
		}
	}

	if v := os.Getenv("TALOS_NODES"); v != "" {
		for _, n := range strings.Split(v, ",") {
			n = strings.TrimSpace(n)
			if n != "" {
				cfg.Nodes = append(cfg.Nodes, n)
			}
		}
	}

	cfg.Perms = loadPermissions()
	return cfg, nil
}

// Permissions controls which categories of operations the MCP server may perform.
// Set MCP_READ_ONLY=true to disable all writes/deletes/actions.
// Fine-grained overrides: MCP_ALLOW_WRITE, MCP_ALLOW_DELETE, MCP_ALLOW_EXEC (values: true/false).
type Permissions struct {
	AllowWrite  bool
	AllowDelete bool
	AllowExec   bool
}

func loadPermissions() Permissions {
	readOnly := os.Getenv("MCP_READ_ONLY") == "true"
	p := Permissions{
		AllowWrite:  !readOnly,
		AllowDelete: !readOnly,
		AllowExec:   !readOnly,
	}
	if v := os.Getenv("MCP_ALLOW_WRITE"); v != "" {
		p.AllowWrite = v == "true"
	}
	if v := os.Getenv("MCP_ALLOW_DELETE"); v != "" {
		p.AllowDelete = v == "true"
	}
	if v := os.Getenv("MCP_ALLOW_EXEC"); v != "" {
		p.AllowExec = v == "true"
	}
	return p
}
