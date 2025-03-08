package main

import (
	"flag"
	"log"
	"os"

	"github.com/ucloud/ucloud-mcp-server/pkg/config"
	"github.com/ucloud/ucloud-mcp-server/pkg/mcp"
	"github.com/ucloud/ucloud-mcp-server/pkg/ucloud"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting UCloud MCP Server...")

	// Define command line flags
	configPath := flag.String("config", "config.json", "Path to configuration file")
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	// Print startup information
	log.Printf("Using config file: %s", *configPath)
	log.Printf("Server will listen on port: %s", *port)

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Failed to load config from file: %v, trying environment variables", err)
		// Try loading from environment variables
		cfg = &config.Config{
			Region:     os.Getenv("UCLOUD_REGION"),
			ProjectID:  os.Getenv("UCLOUD_PROJECT_ID"),
			PublicKey:  os.Getenv("UCLOUD_PUBLIC_KEY"),
			PrivateKey: os.Getenv("UCLOUD_PRIVATE_KEY"),
		}
	}

	// Print configuration (Note: avoid printing sensitive information in production)
	log.Printf("Using configuration - Region: %s, ProjectID: %s", cfg.Region, cfg.ProjectID)

	// Create UCloud client
	ucloudClient, err := ucloud.NewUCloudClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create UCloud client: %v", err)
	}

	// Create MCP server
	mcpServer := mcp.NewMCPServer(ucloudClient)

	// Start server
	if err := mcpServer.Start(*port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
