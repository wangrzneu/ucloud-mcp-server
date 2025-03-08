package utils

import (
	"fmt"
	"strings"
)

// ParsePath parses a path containing variables and returns variable names and their values
// Example: ParsePath("/instances/{instance_id}/status", "/instances/uhost-123/status")
// Returns: map[string]string{"instance_id": "uhost-123"}
func ParsePath(pattern, path string) (map[string]string, error) {
	// Split the path
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Check if path segments match
	if len(patternParts) != len(pathParts) {
		return nil, fmt.Errorf("path length mismatch: pattern has %d parts, path has %d parts",
			len(patternParts), len(pathParts))
	}

	// Store variable names and values
	variables := make(map[string]string)

	// Iterate through each path segment
	for i := 0; i < len(patternParts); i++ {
		patternPart := patternParts[i]
		pathPart := pathParts[i]

		// Check if it's a variable (starts with { and ends with })
		if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			// Extract variable name (remove { and })
			varName := strings.TrimSuffix(strings.TrimPrefix(patternPart, "{"), "}")
			variables[varName] = pathPart
		} else if patternPart != pathPart {
			// If not a variable and doesn't match, return error
			return nil, fmt.Errorf("path segment mismatch at position %d: expected %s, got %s",
				i, patternPart, pathPart)
		}
	}

	return variables, nil
}
