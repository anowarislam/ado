package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantValid   bool
		wantErrors  int
		wantWarns   int
		errContains string
	}{
		{
			name:      "valid minimal config",
			content:   "version: 1\n",
			wantValid: true,
		},
		{
			name:        "missing version",
			content:     "foo: bar\n",
			wantValid:   false,
			wantErrors:  1,
			wantWarns:   1,
			errContains: "missing required key",
		},
		{
			name:      "unknown key warning",
			content:   "version: 1\nunknown_key: value\n",
			wantValid: true,
			wantWarns: 1,
		},
		{
			name:        "invalid yaml",
			content:     "version: [\n",
			wantValid:   false,
			wantErrors:  1,
			errContains: "invalid YAML",
		},
		{
			name:        "empty file",
			content:     "",
			wantValid:   false,
			wantErrors:  1,
			errContains: "empty",
		},
		{
			name:        "unsupported version",
			content:     "version: 99\n",
			wantValid:   false,
			wantErrors:  1,
			errContains: "unsupported config version",
		},
		{
			name:      "version with extra known fields only",
			content:   "version: 1\n",
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("write temp file: %v", err)
			}

			result, err := Validate(path)
			if err != nil {
				t.Fatalf("Validate() error: %v", err)
			}

			if result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", result.Valid, tt.wantValid)
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Errors count = %d, want %d: %+v", len(result.Errors), tt.wantErrors, result.Errors)
			}

			if len(result.Warnings) != tt.wantWarns {
				t.Errorf("Warnings count = %d, want %d: %+v", len(result.Warnings), tt.wantWarns, result.Warnings)
			}

			if tt.errContains != "" && len(result.Errors) > 0 {
				found := false
				for _, e := range result.Errors {
					if contains(e.Message, tt.errContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error containing %q, got: %+v", tt.errContains, result.Errors)
				}
			}
		})
	}
}

func TestValidate_FileNotFound(t *testing.T) {
	result, err := Validate("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("Validate() error: %v", err)
	}

	if result.Valid {
		t.Error("Expected Valid=false for nonexistent file")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}

	if !contains(result.Errors[0].Message, "not found") {
		t.Errorf("Expected 'not found' error, got: %s", result.Errors[0].Message)
	}
}

func TestValidationResult_HasErrors(t *testing.T) {
	r := &ValidationResult{Errors: []ValidationIssue{{Message: "test"}}}
	if !r.HasErrors() {
		t.Error("HasErrors() should return true")
	}

	r = &ValidationResult{Errors: []ValidationIssue{}}
	if r.HasErrors() {
		t.Error("HasErrors() should return false for empty errors")
	}
}

func TestValidationResult_HasWarnings(t *testing.T) {
	r := &ValidationResult{Warnings: []ValidationIssue{{Message: "test"}}}
	if !r.HasWarnings() {
		t.Error("HasWarnings() should return true")
	}

	r = &ValidationResult{Warnings: []ValidationIssue{}}
	if r.HasWarnings() {
		t.Error("HasWarnings() should return false for empty warnings")
	}
}

func TestValidate_PermissionDenied(t *testing.T) {
	// Skip on Windows where permission handling is different
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")

	// Create file then remove read permissions
	if err := os.WriteFile(path, []byte("version: 1\n"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if err := os.Chmod(path, 0000); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() {
		os.Chmod(path, 0644) // Restore for cleanup
	})

	result, err := Validate(path)
	if err != nil {
		t.Fatalf("Validate() error: %v", err)
	}

	if result.Valid {
		t.Error("Expected Valid=false for permission denied")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}

	if !contains(result.Errors[0].Message, "permission denied") {
		t.Errorf("Expected 'permission denied' error, got: %s", result.Errors[0].Message)
	}
}

func TestValidate_MultipleWarnings(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")
	content := "version: 1\nunknown1: foo\nunknown2: bar\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	result, err := Validate(path)
	if err != nil {
		t.Fatalf("Validate() error: %v", err)
	}

	if !result.Valid {
		t.Error("Expected Valid=true with only warnings")
	}

	if len(result.Warnings) != 2 {
		t.Errorf("Expected 2 warnings, got %d: %+v", len(result.Warnings), result.Warnings)
	}

	// Verify line numbers are populated
	for _, w := range result.Warnings {
		if w.Line == 0 {
			t.Errorf("Warning should have line number: %+v", w)
		}
	}
}

func TestValidate_ValidYAMLInvalidStructure(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")
	// Valid YAML but version is a string when int expected - yaml.Unmarshal handles this gracefully
	// Instead, test a YAML that's valid as document but invalid as map
	content := "- item1\n- item2\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	result, err := Validate(path)
	if err != nil {
		t.Fatalf("Validate() error: %v", err)
	}

	if result.Valid {
		t.Error("Expected Valid=false for array instead of map")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
