package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/ucloud/ucloud-mcp-server/pkg/ucloud"
	"github.com/ucloud/ucloud-mcp-server/pkg/utils"
)

// Handlers contains MCP handlers
type Handlers struct {
	ucloudClient *ucloud.UCloudClient
}

// NewHandlers creates new MCP handlers
func NewHandlers(ucloudClient *ucloud.UCloudClient) *Handlers {
	return &Handlers{
		ucloudClient: ucloudClient,
	}
}

// DescribeInstanceHandler handles instance description requests
func (h *Handlers) DescribeInstanceHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceID := request.Params.Arguments["instance_id"].(string)

	instance, err := h.ucloudClient.DescribeInstance(instanceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to describe instance %v: %v", instanceID, err)), nil
	}

	info := ucloud.FormatInstanceInfo(instance)
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal instance info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// GetInstanceMetricsHandler handles instance metrics retrieval requests
func (h *Handlers) GetInstanceMetricsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceID := request.Params.Arguments["instance_id"].(string)

	instance, err := h.ucloudClient.DescribeInstance(instanceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get instance %v: %v", instanceID, err)), nil
	}

	metrics, err := h.ucloudClient.GetInstanceMetrics(instance)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get metrics: %v", err)), nil
	}

	if len(metrics) == 0 {
		return mcp.NewToolResultError("Instance metrics not found"), nil
	}

	// Build metrics response
	metricsResponse := map[string]interface{}{
		"instance_id": instanceID,
		"name":        instance.Name,
		"status":      instance.State,
		"basic_info": map[string]interface{}{
			"cpu":       instance.CPU,
			"memory":    instance.Memory,
			"disk_size": instance.DiskSet[0].Size,
			"zone":      instance.Zone,
			"ip":        instance.IPSet[0].IP,
		},
		"metrics":   metrics,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.MarshalIndent(metricsResponse, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal metrics info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// InstanceStatusHandler handles instance status retrieval requests
func (h *Handlers) InstanceStatusHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	instanceStatusURI := "uhost://instances/{instance_id}/status"
	variables, err := utils.ParsePath(instanceStatusURI, request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path: %v", err)
	}

	instanceID := variables["instance_id"]
	if instanceID == "" {
		return nil, fmt.Errorf("instance_id not found in path")
	}

	instance, err := h.ucloudClient.DescribeInstance(instanceID)
	if err != nil {
		return nil, err
	}

	status := map[string]string{
		"status": instance.State,
	}

	jsonData, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

// InstanceListHandler handles instance list retrieval requests
func (h *Handlers) InstanceListHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	log.Printf("Listing UCloud instances...")

	instances, err := h.ucloudClient.ListInstances()
	if err != nil {
		log.Printf("Error listing instances: %v", err)
		return nil, fmt.Errorf("failed to list instances: %v", err)
	}

	var allInstances []interface{}
	for _, instance := range instances {
		instanceCopy := instance // Create a copy to avoid using loop variable reference
		info := ucloud.FormatInstanceInfo(&instanceCopy)
		allInstances = append(allInstances, info)
		log.Printf("Found instance: ID=%s, Name=%s, Status=%s",
			instance.UHostId, instance.Name, instance.State)
	}

	log.Printf("Total instances found: %d", len(allInstances))

	jsonData, err := json.MarshalIndent(allInstances, "", "  ")
	if err != nil {
		log.Printf("Error marshaling instance data: %v", err)
		return nil, fmt.Errorf("failed to marshal instance data: %v", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

// InstanceListToolHandler handles instance list tool requests
func (h *Handlers) InstanceListToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Listing all UCloud instances...")

	instances, err := h.ucloudClient.ListInstances()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list instances: %v", err)), nil
	}

	var allInstancesWithMetrics []interface{}
	for _, instance := range instances {
		instanceCopy := instance // Create a copy to avoid using loop variable reference
		log.Printf("Processing instance: %s (%s)", instanceCopy.Name, instanceCopy.UHostId)

		// Get instance metrics
		metrics, err := h.ucloudClient.GetInstanceMetrics(&instanceCopy)
		if err != nil {
			log.Printf("Warning: Failed to get metrics for instance %s: %v", instanceCopy.UHostId, err)
		}

		info := ucloud.FormatInstanceInfoWithMetrics(&instanceCopy, metrics)
		allInstancesWithMetrics = append(allInstancesWithMetrics, info)
	}

	log.Printf("Total instances with metrics found: %d", len(allInstancesWithMetrics))

	// Convert to JSON and return
	jsonData, err := json.MarshalIndent(allInstancesWithMetrics, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// InstanceStatusToolHandler handles instance status tool requests
func (h *Handlers) InstanceStatusToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Getting instance status...")

	// No parameters needed, return status for all instances

	instances, err := h.ucloudClient.ListInstances()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list instances: %v", err)), nil
	}

	var statusList []map[string]string
	for _, instance := range instances {
		status := map[string]string{
			"id":     instance.UHostId,
			"name":   instance.Name,
			"status": instance.State,
		}
		statusList = append(statusList, status)
	}

	jsonData, err := json.MarshalIndent(statusList, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal status data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
