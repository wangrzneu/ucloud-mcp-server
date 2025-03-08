package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Config stores UCloud configuration information
type Config struct {
	Region     string `json:"region"`
	ProjectID  string `json:"project_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// LoadConfig loads configuration from file
func LoadConfig(filename string) (*Config, error) {
	// Read file content
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validate required fields
	if config.Region == "" {
		config.Region = os.Getenv("UCLOUD_REGION") // Try reading from environment variables
	}
	if config.ProjectID == "" {
		config.ProjectID = os.Getenv("UCLOUD_PROJECT_ID")
	}
	if config.PublicKey == "" {
		config.PublicKey = os.Getenv("UCLOUD_PUBLIC_KEY")
	}
	if config.PrivateKey == "" {
		config.PrivateKey = os.Getenv("UCLOUD_PRIVATE_KEY")
	}

	// Check if required fields exist
	var missingFields []string
	if config.Region == "" {
		missingFields = append(missingFields, "region")
	}
	if config.ProjectID == "" {
		missingFields = append(missingFields, "project_id")
	}
	if config.PublicKey == "" {
		missingFields = append(missingFields, "public_key")
	}
	if config.PrivateKey == "" {
		missingFields = append(missingFields, "private_key")
	}

	if len(missingFields) > 0 {
		return nil, fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	return &config, nil
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	return &Config{
		Region:     os.Getenv("UCLOUD_REGION"),
		ProjectID:  os.Getenv("UCLOUD_PROJECT_ID"),
		PublicKey:  os.Getenv("UCLOUD_PUBLIC_KEY"),
		PrivateKey: os.Getenv("UCLOUD_PRIVATE_KEY"),
	}
}
