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

	return cfg, nil
}
