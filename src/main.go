package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"

	"github.com/immxmmi/talos-mcp/config"
	"github.com/immxmmi/talos-mcp/talos"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	client, err := talos.NewClient(cfg)
	if err != nil {
		log.Fatalf("talos client error: %v", err)
	}
	defer client.Close()

	s := server.NewMCPServer("talos-mcp", "1.0.0")
	talos.RegisterTools(s, client)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
