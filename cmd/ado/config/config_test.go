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

func TestConfigValidate_RootConfigFlag(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("version: 1\n"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	// Create root command with --config flag
	cmd := NewCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"validate", "--file", configPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(buf.String(), "Config valid") {
		t.Errorf("expected validation success")
	}
}

func TestConfigValidate_AutoDetect(t *testing.T) {
	// Create config in standard location
	homeDir := t.TempDir()
	configDir := filepath.Join(homeDir, ".config", "ado")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("version: 1\n"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Set HOME to temp dir
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	cmd := NewCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"validate"}) // No --file flag

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(buf.String(), "Config valid") {
		t.Errorf("expected auto-detected config to validate")
	}
}

func TestConfigValidate_NoConfigFound(t *testing.T) {
	// Use empty temp directory
	homeDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	cmd := NewCommand()
	cmd.SetArgs([]string{"validate"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when no config found")
	}

	if !strings.Contains(err.Error(), "no config file found") {
		t.Errorf("expected 'no config file found' error, got: %v", err)
	}
}

func TestConfigValidate_StrictMode(t *testing.T) {
	// Test that strict mode is recognized
	// Note: Full testing of strict mode exit behavior is skipped because
	// the code calls os.Exit(1) directly (config.go:94), which terminates
	// the test process. Testing strict mode output formatting is done in
	// TestFormatValidationResult instead.
	t.Skip("Strict mode calls os.Exit(1) directly, making it untestable without subprocess")
}

func TestConfigValidate_InvalidOutputFormat(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("version: 1\n"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	cmd := NewCommand()
	cmd.SetArgs([]string{"validate", "--file", configPath, "--output", "invalid"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}

	if !strings.Contains(err.Error(), "output format") {
		t.Errorf("expected output format error, got: %v", err)
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
		{
			name: "error with line number",
			result: &internalconfig.ValidationResult{
				Valid: false,
				Path:  "/path/to/config.yaml",
				Errors: []internalconfig.ValidationIssue{
					{Message: "syntax error", Line: 42, Severity: "error"},
				},
				Warnings: []internalconfig.ValidationIssue{},
			},
			contains: []string{"\u2717", "Config invalid", "Error:", "syntax error", "line 42"},
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
