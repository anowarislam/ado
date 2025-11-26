package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ValidationResult holds the result of config validation.
type ValidationResult struct {
	Valid    bool              `json:"valid" yaml:"valid"`
	Path     string            `json:"path" yaml:"path"`
	Errors   []ValidationIssue `json:"errors" yaml:"errors"`
	Warnings []ValidationIssue `json:"warnings" yaml:"warnings"`
}

// ValidationIssue represents a single validation error or warning.
type ValidationIssue struct {
	Message  string `json:"message" yaml:"message"`
	Line     int    `json:"line,omitempty" yaml:"line,omitempty"`
	Severity string `json:"severity" yaml:"severity"`
}

// ConfigSchema represents the expected config file structure.
type ConfigSchema struct {
	Version int `yaml:"version"`
}

// knownKeys lists valid top-level config keys.
var knownKeys = map[string]bool{
	"version": true,
}

// Validate validates a config file at the given path.
// Returns a ValidationResult with any errors or warnings found.
func Validate(path string) (*ValidationResult, error) {
	result := &ValidationResult{
		Path:     path,
		Valid:    true,
		Errors:   []ValidationIssue{},
		Warnings: []ValidationIssue{},
	}

	// Check file exists
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationIssue{
				Message:  fmt.Sprintf("config file not found: %q", path),
				Severity: "error",
			})
			return result, nil
		}
		if os.IsPermission(err) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationIssue{
				Message:  fmt.Sprintf("permission denied: %q", path),
				Severity: "error",
			})
			return result, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Message:  "config file is empty",
			Severity: "error",
		})
		return result, nil
	}

	// Parse YAML to check syntax and get line numbers
	var rawNode yaml.Node
	if err := yaml.Unmarshal(data, &rawNode); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Message:  fmt.Sprintf("invalid YAML: %s", err.Error()),
			Severity: "error",
		})
		return result, nil
	}

	// Parse into map to check for unknown keys
	var rawMap map[string]any
	if err := yaml.Unmarshal(data, &rawMap); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Message:  fmt.Sprintf("invalid YAML structure: %s", err.Error()),
			Severity: "error",
		})
		return result, nil
	}

	// Check for unknown keys
	for key := range rawMap {
		if !knownKeys[key] {
			line := findKeyLine(&rawNode, key)
			result.Warnings = append(result.Warnings, ValidationIssue{
				Message:  fmt.Sprintf("unknown key %q", key),
				Line:     line,
				Severity: "warning",
			})
		}
	}

	// Parse into schema struct for validation
	var schema ConfigSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Message:  fmt.Sprintf("invalid config structure: %s", err.Error()),
			Severity: "error",
		})
		return result, nil
	}

	// Validate required fields
	if schema.Version == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Message:  "missing required key \"version\"",
			Severity: "error",
		})
	} else if schema.Version != 1 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Message:  fmt.Sprintf("unsupported config version: %d (expected: 1)", schema.Version),
			Severity: "error",
		})
	}

	return result, nil
}

// findKeyLine searches the YAML node tree for a key and returns its line number.
func findKeyLine(node *yaml.Node, key string) int {
	if node == nil {
		return 0
	}

	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			if line := findKeyLine(child, key); line > 0 {
				return line
			}
		}
	case yaml.MappingNode:
		for i := 0; i < len(node.Content)-1; i += 2 {
			keyNode := node.Content[i]
			if keyNode.Value == key {
				return keyNode.Line
			}
		}
	}
	return 0
}

// HasErrors returns true if there are any validation errors.
func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// HasWarnings returns true if there are any validation warnings.
func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}
