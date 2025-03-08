package ucloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ucloud/ucloud-mcp-server/pkg/config"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

// UCloudClient wraps the UCloud API client
type UCloudClient struct {
	UHostClient   *uhost.UHostClient
	GenericClient *ucloud.Client
}

// NewUCloudClient creates a new UCloud client
func NewUCloudClient(cfg *config.Config) (*UCloudClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration is nil")
	}

	// Create UCloud configuration
	ucfg := ucloud.NewConfig()
	ucfg.Region = cfg.Region
	ucfg.ProjectId = cfg.ProjectID
	ucfg.BaseUrl = "https://api.ucloud.cn"

	// Create credentials
	credential := auth.NewCredential()
	credential.PublicKey = cfg.PublicKey
	credential.PrivateKey = cfg.PrivateKey

	// Create UHost client
	uhostClient := uhost.NewClient(&ucfg, &credential)

	// Create generic client
	genericClient := ucloud.NewClient(&ucfg, &credential)

	return &UCloudClient{
		UHostClient:   uhostClient,
		GenericClient: genericClient,
	}, nil
}

// DescribeInstance gets detailed information about an instance
func (c *UCloudClient) DescribeInstance(instanceID string) (*uhost.UHostInstanceSet, error) {
	req := c.UHostClient.NewDescribeUHostInstanceRequest()
	req.UHostIds = []string{instanceID}

	resp, err := c.UHostClient.DescribeUHostInstance(req)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance %s: %v", instanceID, err)
	}

	if len(resp.UHostSet) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}

	return &resp.UHostSet[0], nil
}

// ListInstances gets a list of instances
func (c *UCloudClient) ListInstances() ([]uhost.UHostInstanceSet, error) {
	var allInstances []uhost.UHostInstanceSet
	limit := 100
	offset := 0

	for {
		req := c.UHostClient.NewDescribeUHostInstanceRequest()
		req.Limit = &limit
		req.Offset = &offset

		resp, err := c.UHostClient.DescribeUHostInstance(req)
		if err != nil {
			return nil, fmt.Errorf("failed to list instances: %v", err)
		}

		allInstances = append(allInstances, resp.UHostSet...)

		// If the number of instances is less than limit, we've got all data
		if len(resp.UHostSet) < limit {
			break
		}

		// Update offset for next page
		offset += limit
	}

	return allInstances, nil
}

// InstanceMetrics represents instance monitoring metrics
type InstanceMetrics struct {
	ResourceId     string  `json:"ResourceId"`
	CPUUtilization float64 `json:"CPUUtilization"`
	IORead         float64 `json:"IORead"`
	IOWrite        float64 `json:"IOWrite"`
	DiskReadOps    float64 `json:"DiskReadOps"`
	DiskWriteOps   float64 `json:"DiskWriteOps"`
	NICIn          float64 `json:"NICIn"`
	NICOut         float64 `json:"NICOut"`
	NetPacketIn    float64 `json:"NetPacketIn"`
	NetPacketOut   float64 `json:"NetPacketOut"`
	Name           string  `json:"Name"`
	PrivateIp      string  `json:"PrivateIp"`
	CreateTime     int64   `json:"CreateTime"`
}

// GetInstanceMetrics retrieves instance monitoring metrics
func (c *UCloudClient) GetInstanceMetrics(instance *uhost.UHostInstanceSet) ([]InstanceMetrics, error) {
	if instance == nil {
		return nil, fmt.Errorf("instance is nil")
	}

	var metrics []InstanceMetrics
	limit := 100
	offset := 0

	for {
		req := c.GenericClient.NewGenericRequest()
		err := req.SetPayload(map[string]interface{}{
			"Action":       "GetMetricOverview",
			"Zone":         instance.Zone,
			"ResourceType": "uhost",
			"Limit":        limit,
			"Offset":       offset,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to set payload: %v", err)
		}

		// Call monitoring API
		metricResp, err := c.GenericClient.GenericInvoke(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get metrics: %v", err)
		}

		// Get DataSet from response
		var metricsData struct {
			RetCode      int               `json:"RetCode"`
			Action       string            `json:"Action"`
			ResourceType string            `json:"ResourceType"`
			DataSet      []InstanceMetrics `json:"DataSet"`
			RefreshTime  int64             `json:"RefreshTime"`
			TotalCount   int               `json:"TotalCount"`
		}

		payload := metricResp.GetPayload()
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %v", err)
		}

		if err := json.Unmarshal(jsonBytes, &metricsData); err != nil {
			return nil, fmt.Errorf("failed to parse metrics data: %v", err)
		}

		// Find monitoring data for the specified instance
		for _, data := range metricsData.DataSet {
			if data.ResourceId == instance.UHostId {
				metrics = append(metrics, data)
			}
		}

		// If the number of data points is less than limit, we've got all data
		if len(metricsData.DataSet) < limit {
			break
		}

		// Update offset for next page
		offset += limit

		// If we've reached TotalCount, we can exit
		if offset >= metricsData.TotalCount {
			break
		}
	}

	return metrics, nil
}

// InstanceInfo represents instance information for API response
type InstanceInfo struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	IP        string      `json:"ip"`
	Zone      string      `json:"zone"`
	CPU       int         `json:"cpu"`
	Memory    int         `json:"memory"`
	DiskSize  int         `json:"disk_size"`
	Metrics   interface{} `json:"metrics,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

// FormatInstanceInfo formats UHost instance information for API response
func FormatInstanceInfo(instance *uhost.UHostInstanceSet) *InstanceInfo {
	if instance == nil {
		return nil
	}

	var ip string
	if len(instance.IPSet) > 0 {
		ip = instance.IPSet[0].IP
	}

	var diskSize int
	if len(instance.DiskSet) > 0 {
		diskSize = instance.DiskSet[0].Size
	}

	return &InstanceInfo{
		ID:        instance.UHostId,
		Name:      instance.Name,
		Status:    instance.State,
		IP:        ip,
		Zone:      instance.Zone,
		CPU:       instance.CPU,
		Memory:    instance.Memory,
		DiskSize:  diskSize,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// FormatInstanceInfoWithMetrics formats UHost instance information with monitoring data for API response
func FormatInstanceInfoWithMetrics(instance *uhost.UHostInstanceSet, metrics []InstanceMetrics) *InstanceInfo {
	info := FormatInstanceInfo(instance)
	if info != nil {
		info.Metrics = metrics
	}
	return info
}
