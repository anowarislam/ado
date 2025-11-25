package meta

import (
	"bytes"
	"strings"
	"testing"

	internalmeta "github.com/anowarislam/ado/internal/meta"
)

func TestNewCommand(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{
		Name:      "ado",
		Version:   "1.0.0",
		Commit:    "abc123",
		BuildTime: "2024-01-01T00:00:00Z",
		GoVersion: "go1.22.0",
		Platform:  "darwin/arm64",
	}

	cmd := NewCommand(buildInfo)

	if cmd.Use != "meta" {
		t.Errorf("Use = %q, want %q", cmd.Use, "meta")
	}

	// Verify subcommands
	subcommands := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subcommands[sub.Name()] = true
	}

	expectedCmds := []string{"info", "env", "features"}
	for _, name := range expectedCmds {
		if !subcommands[name] {
			t.Errorf("expected subcommand %q not found", name)
		}
	}
}

func TestMetaInfo(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{
		Name:      "ado",
		Version:   "1.0.0",
		Commit:    "abc123",
		BuildTime: "2024-01-01T00:00:00Z",
		GoVersion: "go1.22.0",
		Platform:  "darwin/arm64",
	}

	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"info"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedFields := []string{"Name:", "Version:", "Commit:", "BuildTime:", "GoVersion:", "Platform:"}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("output missing %q", field)
		}
	}
}

func TestMetaInfo_JSON(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{
		Name:    "ado",
		Version: "1.0.0",
	}

	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"info", "--output", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"name"`) || !strings.Contains(output, `"version"`) {
		t.Errorf("JSON output missing expected fields: %s", output)
	}
}

func TestMetaFeatures(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"features"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No experimental features") {
		t.Errorf("output = %q, expected 'No experimental features'", output)
	}
}

func TestFormatBuildInfo(t *testing.T) {
	info := internalmeta.BuildInfo{
		Name:      "ado",
		Version:   "1.0.0",
		Commit:    "abc123",
		BuildTime: "2024-01-01",
		GoVersion: "go1.22.0",
		Platform:  "linux/amd64",
	}

	output := formatBuildInfo(info)

	if !strings.Contains(output, "Name: ado") {
		t.Error("missing Name field")
	}
	if !strings.Contains(output, "Version: 1.0.0") {
		t.Error("missing Version field")
	}
}

func TestFormatEnvInfo(t *testing.T) {
	info := internalmeta.EnvInfo{
		ConfigPath:    "/path/to/config",
		ConfigSources: []string{"/source1", "/source2"},
		HomeDir:       "/home/user",
		CacheDir:      "/cache",
		Env:           map[string]string{"FOO": "bar"},
	}

	output := formatEnvInfo(info)

	if !strings.Contains(output, "ConfigPath: /path/to/config") {
		t.Error("missing ConfigPath")
	}
	if !strings.Contains(output, "/source1") {
		t.Error("missing ConfigSources")
	}
	if !strings.Contains(output, "FOO=bar") {
		t.Error("missing EnvVariables")
	}
}

func TestFormatEnvInfo_Empty(t *testing.T) {
	info := internalmeta.EnvInfo{
		ConfigPath:    "",
		ConfigSources: []string{},
		HomeDir:       "/home",
		CacheDir:      "/cache",
		Env:           map[string]string{},
	}

	output := formatEnvInfo(info)

	if !strings.Contains(output, "(none resolved)") {
		t.Error("missing '(none resolved)' for empty config path")
	}
	if !strings.Contains(output, "(none set)") {
		t.Error("missing '(none set)' for empty env")
	}
}
