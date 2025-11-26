package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name:   "default config",
			config: DefaultConfig(),
		},
		{
			name: "debug level",
			config: Config{
				Level:  "debug",
				Format: "text",
				Output: "stderr",
			},
		},
		{
			name: "json format",
			config: Config{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.config)
			if logger == nil {
				t.Fatal("New() returned nil")
			}
			if logger.Handler() == nil {
				t.Error("Handler() returned nil")
			}
		})
	}
}

func TestDefault(t *testing.T) {
	logger := Default()
	if logger == nil {
		t.Fatal("Default() returned nil")
	}
}

func TestNopLogger(t *testing.T) {
	logger := NopLogger()
	if logger == nil {
		t.Fatal("NopLogger() returned nil")
	}

	// Should not panic
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
}

func TestLoggerWith(t *testing.T) {
	logger := Default()
	childLogger := logger.With("key", "value")

	if childLogger == nil {
		t.Fatal("With() returned nil")
	}

	// Should be a new instance
	if logger == childLogger {
		t.Error("With() should return a new logger instance")
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name         string
		configLevel  string
		logLevel     string
		shouldAppear bool
	}{
		{"debug at debug level", "debug", "debug", true},
		{"info at debug level", "debug", "info", true},
		{"warn at debug level", "debug", "warn", true},
		{"error at debug level", "debug", "error", true},

		{"debug at info level", "info", "debug", false},
		{"info at info level", "info", "info", true},
		{"warn at info level", "info", "warn", true},
		{"error at info level", "info", "error", true},

		{"debug at warn level", "warn", "debug", false},
		{"info at warn level", "warn", "info", false},
		{"warn at warn level", "warn", "warn", true},
		{"error at warn level", "warn", "error", true},

		{"debug at error level", "error", "debug", false},
		{"info at error level", "error", "info", false},
		{"warn at error level", "error", "warn", false},
		{"error at error level", "error", "error", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newTestLogger(tt.configLevel, &buf)

			// Log at the specified level
			switch tt.logLevel {
			case "debug":
				logger.Debug("test message")
			case "info":
				logger.Info("test message")
			case "warn":
				logger.Warn("test message")
			case "error":
				logger.Error("test message")
			}

			output := buf.String()
			hasOutput := len(output) > 0

			if hasOutput != tt.shouldAppear {
				t.Errorf("expected shouldAppear=%v, got output=%q", tt.shouldAppear, output)
			}
		})
	}
}

func TestJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLoggerWithFormat("info", "json", &buf)

	logger.Info("test message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Fatal("expected output, got empty string")
	}

	// Should be valid JSON
	var data map[string]any
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Errorf("output is not valid JSON: %v\nOutput: %s", err, output)
	}

	// Should contain expected fields
	if data["msg"] != "test message" {
		t.Errorf("expected msg='test message', got %v", data["msg"])
	}
	if data["key"] != "value" {
		t.Errorf("expected key='value', got %v", data["key"])
	}
}

func TestTextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLoggerWithFormat("info", "text", &buf)

	logger.Info("test message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Fatal("expected output, got empty string")
	}

	// Should contain message and key-value
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain 'test message': %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("output should contain 'key=value': %s", output)
	}
}

func TestStructuredFields(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLoggerWithFormat("info", "json", &buf)

	logger.Info("operation complete",
		"user_id", 123,
		"action", "login",
		"success", true,
	)

	output := buf.String()

	var data map[string]any
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Check numeric field
	if data["user_id"] != float64(123) {
		t.Errorf("expected user_id=123, got %v", data["user_id"])
	}

	// Check string field
	if data["action"] != "login" {
		t.Errorf("expected action='login', got %v", data["action"])
	}

	// Check boolean field
	if data["success"] != true {
		t.Errorf("expected success=true, got %v", data["success"])
	}
}

func TestWithContext(t *testing.T) {
	logger := Default()
	ctx := context.Background()

	// No logger in context should return default
	fromCtx := FromContext(ctx)
	if fromCtx == nil {
		t.Fatal("FromContext() returned nil for context without logger")
	}

	// Add logger to context
	ctx = WithContext(ctx, logger)
	fromCtx = FromContext(ctx)
	if fromCtx == nil {
		t.Fatal("FromContext() returned nil for context with logger")
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"invalid", slog.LevelInfo},
		{"", slog.LevelInfo},
		{"DEBUG", slog.LevelInfo}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestResolveFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		isTTY    bool
		expected string
	}{
		{"json explicit", "json", true, "json"},
		{"json explicit non-tty", "json", false, "json"},
		{"text explicit", "text", true, "text"},
		{"text explicit non-tty", "text", false, "text"},
		{"auto with tty", "auto", true, "text"},
		{"auto without tty", "auto", false, "json"},
		{"invalid defaults to text", "invalid", true, "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output io.Writer
			if tt.isTTY {
				// Use a mock that reports as TTY - we can't easily mock this
				// so we just test the explicit formats
				if tt.format == "auto" {
					// Skip TTY-dependent test
					return
				}
			}
			output = &bytes.Buffer{}

			result := resolveFormat(tt.format, output)
			if tt.format != "auto" && result != tt.expected {
				t.Errorf("resolveFormat(%q) = %q, want %q", tt.format, result, tt.expected)
			}
		})
	}
}

func TestResolveOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{"stdout", "stdout", "stdout"},
		{"stderr", "stderr", "stderr"},
		{"default to stderr", "invalid", "stderr"},
		{"empty defaults to stderr", "", "stderr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveOutput(tt.output)
			// We can't directly compare io.Writer, but we can check it's not nil
			if result == nil {
				t.Error("resolveOutput() returned nil")
			}
		})
	}
}

func TestNewWithStdout(t *testing.T) {
	// Test creating a logger with stdout output
	logger := New(Config{
		Level:  "info",
		Format: "text",
		Output: "stdout",
	})

	if logger == nil {
		t.Fatal("New() with stdout returned nil")
	}
}

// newTestLogger creates a logger that writes to the provided buffer.
func newTestLogger(level string, buf *bytes.Buffer) Logger {
	return newTestLoggerWithFormat(level, "text", buf)
}

// newTestLoggerWithFormat creates a logger with specific format.
func newTestLoggerWithFormat(level, format string, buf *bytes.Buffer) Logger {
	opts := &slog.HandlerOptions{
		Level: parseLevel(level),
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(buf, opts)
	} else {
		handler = slog.NewTextHandler(buf, opts)
	}

	return &logger{
		slog: slog.New(handler),
	}
}
