package logging

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Level != "info" {
		t.Errorf("Level = %q, want %q", cfg.Level, "info")
	}
	if cfg.Format != "auto" {
		t.Errorf("Format = %q, want %q", cfg.Format, "auto")
	}
	if cfg.Output != "stderr" {
		t.Errorf("Output = %q, want %q", cfg.Output, "stderr")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name     string
		input    Config
		expected Config
	}{
		{
			name:     "valid config unchanged",
			input:    Config{Level: "debug", Format: "json", Output: "stdout"},
			expected: Config{Level: "debug", Format: "json", Output: "stdout"},
		},
		{
			name:     "invalid level defaults to info",
			input:    Config{Level: "invalid", Format: "text", Output: "stderr"},
			expected: Config{Level: "info", Format: "text", Output: "stderr"},
		},
		{
			name:     "invalid format defaults to auto",
			input:    Config{Level: "info", Format: "invalid", Output: "stderr"},
			expected: Config{Level: "info", Format: "auto", Output: "stderr"},
		},
		{
			name:     "invalid output defaults to stderr",
			input:    Config{Level: "info", Format: "text", Output: "invalid"},
			expected: Config{Level: "info", Format: "text", Output: "stderr"},
		},
		{
			name:     "empty config gets defaults",
			input:    Config{},
			expected: Config{Level: "info", Format: "auto", Output: "stderr"},
		},
		{
			name:     "all valid levels",
			input:    Config{Level: "debug", Format: "auto", Output: "stderr"},
			expected: Config{Level: "debug", Format: "auto", Output: "stderr"},
		},
		{
			name:     "warn level",
			input:    Config{Level: "warn", Format: "auto", Output: "stderr"},
			expected: Config{Level: "warn", Format: "auto", Output: "stderr"},
		},
		{
			name:     "error level",
			input:    Config{Level: "error", Format: "auto", Output: "stderr"},
			expected: Config{Level: "error", Format: "auto", Output: "stderr"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Validate()

			if result.Level != tt.expected.Level {
				t.Errorf("Level = %q, want %q", result.Level, tt.expected.Level)
			}
			if result.Format != tt.expected.Format {
				t.Errorf("Format = %q, want %q", result.Format, tt.expected.Format)
			}
			if result.Output != tt.expected.Output {
				t.Errorf("Output = %q, want %q", result.Output, tt.expected.Output)
			}
		})
	}
}

func TestIsValidLevel(t *testing.T) {
	tests := []struct {
		level string
		valid bool
	}{
		{"debug", true},
		{"info", true},
		{"warn", true},
		{"error", true},
		{"DEBUG", false}, // Case sensitive
		{"Info", false},
		{"invalid", false},
		{"", false},
		{"trace", false},
		{"fatal", false},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			result := IsValidLevel(tt.level)
			if result != tt.valid {
				t.Errorf("IsValidLevel(%q) = %v, want %v", tt.level, result, tt.valid)
			}
		})
	}
}
