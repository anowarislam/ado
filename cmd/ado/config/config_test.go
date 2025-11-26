package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	internalconfig "github.com/anowarislam/ado/internal/config"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	if cmd.Use != "config" {
		t.Errorf("Use = %q, want %q", cmd.Use, "config")
	}

	// Verify subcommands
	subcommands := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = true
	}

	if !subcommands["validate"] {
		t.Error("expected subcommand 'validate' not found")
	}
}

func TestConfigValidate_ValidFile(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("version: 1\n"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	cmd := NewCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"validate", "--file", configPath})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Config valid") {
		t.Errorf("expected 'Config valid' in output, got: %s", output)
	}
	if !strings.Contains(output, "\u2713") {
		t.Errorf("expected checkmark in output, got: %s", output)
	}
}

func TestConfigValidate_WithWarning(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("version: 1\nunknown_key: value\n"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	cmd := NewCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"validate", "--file", configPath})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Config valid") {
		t.Errorf("expected 'Config valid' in output, got: %s", output)
	}
	if !strings.Contains(output, "Warning") {
		t.Errorf("expected warning in output, got: %s", output)
	}
	if !strings.Contains(output, "unknown_key") {
		t.Errorf("expected 'unknown_key' in warning, got: %s", output)
	}
}

func TestConfigValidate_JSONOutput(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("version: 1\n"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	cmd := NewCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"validate", "--file", configPath, "--output", "json"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"valid": true`) {
		t.Errorf("expected JSON with valid=true, got: %s", output)
	}
	if !strings.Contains(output, `"path"`) {
		t.Errorf("expected JSON with path field, got: %s", output)
	}
}

func TestFormatValidationResult(t *testing.T) {
	tests := []struct {
		name     string
		result   *internalconfig.ValidationResult
		contains []string
	}{
		{
			name: "valid config",
			result: &internalconfig.ValidationResult{
				Valid:    true,
				Path:     "/path/to/config.yaml",
				Errors:   []internalconfig.ValidationIssue{},
				Warnings: []internalconfig.ValidationIssue{},
			},
			contains: []string{"\u2713", "Config valid", "/path/to/config.yaml"},
		},
		{
			name: "invalid config with error",
			result: &internalconfig.ValidationResult{
				Valid: false,
				Path:  "/path/to/config.yaml",
				Errors: []internalconfig.ValidationIssue{
					{Message: "missing version", Severity: "error"},
				},
				Warnings: []internalconfig.ValidationIssue{},
			},
			contains: []string{"\u2717", "Config invalid", "Error:", "missing version"},
		},
		{
			name: "valid with warning",
			result: &internalconfig.ValidationResult{
				Valid:  true,
				Path:   "/path/to/config.yaml",
				Errors: []internalconfig.ValidationIssue{},
				Warnings: []internalconfig.ValidationIssue{
					{Message: "unknown key", Line: 5, Severity: "warning"},
				},
			},
			contains: []string{"\u2713", "Config valid", "Warning:", "unknown key", "line 5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := formatValidationResult(tt.result)
			for _, substr := range tt.contains {
				if !strings.Contains(output, substr) {
					t.Errorf("output missing %q: %s", substr, output)
				}
			}
		})
	}
}
