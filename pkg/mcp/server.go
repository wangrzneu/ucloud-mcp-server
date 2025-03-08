package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ucloud/ucloud-mcp-server/pkg/ucloud"
)

// MCPServer wraps the MCP server implementation
type MCPServer struct {
	server       *server.MCPServer
	sseServer    *server.SSEServer
	handlers     *Handlers
	ucloudClient *ucloud.UCloudClient
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(ucloudClient *ucloud.UCloudClient) *MCPServer {
	// Create MCP server
	mcpServer := server.NewMCPServer(
		"UCloud Instance Manager",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	handlers := NewHandlers(ucloudClient)

	return &MCPServer{
		server:       mcpServer,
		handlers:     handlers,
		ucloudClient: ucloudClient,
	}
}

// RegisterTools registers all tools
func (s *MCPServer) RegisterTools() {
	// Add describe instance tool
	describeTool := mcp.NewTool("describe_instance",
		mcp.WithDescription("Get information about a UCloud instance"),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("ID of the instance to describe"),
		),
	)
	s.server.AddTool(describeTool, s.handlers.DescribeInstanceHandler)

	// Add monitoring metrics tool
	monitorTool := mcp.NewTool("get_instance_metrics",
		mcp.WithDescription("Get monitoring metrics for a UCloud instance"),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("ID of the instance to monitor"),
		),
	)
	s.server.AddTool(monitorTool, s.handlers.GetInstanceMetricsHandler)

	// Add instance status tool
	instanceStatusTool := mcp.NewTool("instance_status",
		mcp.WithDescription("Get the current status of a UCloud instance"),
		mcp.WithString("random_string",
			mcp.Required(),
			mcp.Description("Dummy parameter for no-parameter tools"),
		),
	)
	s.server.AddTool(instanceStatusTool, s.handlers.InstanceStatusToolHandler)

	// Add instance list tool
	instanceListTool := mcp.NewTool("instance_list",
		mcp.WithDescription("List all UCloud instances"),
		mcp.WithString("random_string",
			mcp.Required(),
			mcp.Description("Dummy parameter for no-parameter tools"),
		),
	)
	s.server.AddTool(instanceListTool, s.handlers.InstanceListToolHandler)
}

// RegisterResources registers all resources
func (s *MCPServer) RegisterResources() {
	// Add instance status resource
	instanceStatusURI := "uhost://instances/{instance_id}/status"
	s.server.AddResource(mcp.NewResource(instanceStatusURI, "instance_status",
		mcp.WithResourceDescription("Get the current status of a UCloud instance"),
		mcp.WithMIMEType("application/json"),
	), s.handlers.InstanceStatusHandler)

	// Add instance list resource
	s.server.AddResource(mcp.NewResource("uhost://instances", "instance_list",
		mcp.WithResourceDescription("List all UCloud instances"),
		mcp.WithMIMEType("application/json"),
	), s.handlers.InstanceListHandler)
}

// RegisterPrompts registers all prompts
func (s *MCPServer) RegisterPrompts() {
	// Add instance management prompt
	s.server.AddPrompt(mcp.NewPrompt("instance_management",
		mcp.WithPromptDescription("Help with UCloud instance management"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("Action to perform (describe, list, get_instance_metrics)"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]

		var messages []mcp.PromptMessage
		switch action {
		case "describe":
			messages = []mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent("I'll help you get information about a UCloud instance."),
				),
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent("Please provide the instance ID you want to describe."),
				),
			}
		case "list":
			messages = []mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent("I'll help you list all your UCloud instances."),
				),
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent("I'll show you a list of all instances with their details."),
				),
			}
		case "get_instance_metrics":
			messages = []mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent("I'll help you get the monitoring metrics for a UCloud instance."),
				),
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent("Please provide the instance ID you want to get metrics for."),
				),
			}
		default:
			return nil, fmt.Errorf("unknown action: %s", action)
		}

		return mcp.NewGetPromptResult(
			"UCloud Instance Management",
			messages,
		), nil
	})
}

// Start starts the SSE server
func (s *MCPServer) Start(port string) error {
	// Register all tools, resources and prompts
	s.RegisterTools()
	s.RegisterResources()
	s.RegisterPrompts()

	// Create and start SSE server
	s.sseServer = server.NewSSEServer(s.server, " http://localhost:"+port)
	log.Printf("SSE server listening on :%s", port)
	return s.sseServer.Start(":" + port)
}
