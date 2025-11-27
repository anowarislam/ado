package meta

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

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

	expectedCmds := []string{"info", "env", "features", "system"}
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

func TestMetaEnv(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)

	// Set up root command context with config flag (required by env command)
	root := &cobra.Command{Use: "ado"}
	root.PersistentFlags().String("config", "", "Path to config file")
	root.AddCommand(cmd)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"meta", "env"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedFields := []string{"ConfigPath:", "ConfigSources:", "HomeDir:", "CacheDir:", "EnvVariables:"}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("output missing %q", field)
		}
	}
}

func TestMetaEnv_JSON(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)

	root := &cobra.Command{Use: "ado"}
	root.PersistentFlags().String("config", "", "Path to config file")
	root.AddCommand(cmd)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"meta", "env", "--output", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"config_path"`) || !strings.Contains(output, `"home_dir"`) {
		t.Errorf("JSON output missing expected fields: %s", output)
	}
}

func TestMetaEnv_YAML(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)

	root := &cobra.Command{Use: "ado"}
	root.PersistentFlags().String("config", "", "Path to config file")
	root.AddCommand(cmd)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"meta", "env", "--output", "yaml"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "config_path:") || !strings.Contains(output, "home_dir:") {
		t.Errorf("YAML output missing expected fields: %s", output)
	}
}

func TestMetaEnv_InvalidOutput(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)

	root := &cobra.Command{Use: "ado"}
	root.PersistentFlags().String("config", "", "Path to config file")
	root.AddCommand(cmd)

	root.SetArgs([]string{"meta", "env", "--output", "invalid"})

	err := root.Execute()
	if err == nil {
		t.Error("expected error for invalid output format")
	}
}

func TestMetaInfo_InvalidOutput(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"info", "--output", "invalid"})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid output format")
	}
}

func TestMetaInfo_YAML(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{
		Name:    "ado",
		Version: "1.0.0",
	}

	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"info", "--output", "yaml"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "name:") || !strings.Contains(output, "version:") {
		t.Errorf("YAML output missing expected fields: %s", output)
	}
}

func TestMetaFeatures_JSON(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"features", "--output", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"features"`) {
		t.Errorf("JSON output missing 'features' field: %s", output)
	}
}

func TestMetaFeatures_YAML(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"features", "--output", "yaml"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "features:") {
		t.Errorf("YAML output missing 'features' field: %s", output)
	}
}

func TestMetaFeatures_InvalidOutput(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)
	cmd.SetArgs([]string{"features", "--output", "invalid"})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid output format")
	}
}

func TestMetaSystem(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"system"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedFields := []string{"OS:", "Platform:", "Kernel:", "Architecture:", "CPU:", "Memory:"}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("output missing %q", field)
		}
	}
}

func TestMetaSystem_JSON(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"system", "--output", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedFields := []string{`"os"`, `"platform"`, `"cpu"`, `"memory"`, `"storage"`, `"gpu"`}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("JSON output missing %q", field)
		}
	}
}

func TestMetaSystem_YAML(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"system", "--output", "yaml"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedFields := []string{"os:", "platform:", "cpu:", "memory:", "storage:", "gpu:"}
	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("YAML output missing %q", field)
		}
	}
}

func TestMetaSystem_InvalidOutput(t *testing.T) {
	buildInfo := internalmeta.BuildInfo{}
	cmd := NewCommand(buildInfo)
	cmd.SetArgs([]string{"system", "--output", "invalid"})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid output format")
	}
}

func TestFormatSystemInfo(t *testing.T) {
	info := internalmeta.SystemInfo{
		OS:           "darwin",
		Platform:     "macOS 14.2",
		Kernel:       "Darwin 23.2.0",
		Architecture: "arm64",
		CPU: internalmeta.CPUInfo{
			Model:        "Apple M2 Pro",
			Vendor:       "Apple",
			Cores:        10,
			FrequencyMHz: 0,
		},
		Memory: internalmeta.MemoryInfo{
			TotalMB:     16384,
			AvailableMB: 8192,
			UsedMB:      8192,
			UsedPercent: 50.0,
			SwapTotalMB: 0,
			SwapUsedMB:  0,
		},
		Storage: []internalmeta.StorageInfo{
			{
				Device:      "/dev/disk3s1s1",
				Mountpoint:  "/",
				Filesystem:  "apfs",
				TotalMB:     505856,
				UsedMB:      125952,
				FreeMB:      379904,
				UsedPercent: 25.0,
			},
		},
		GPU: []internalmeta.GPUInfo{
			{
				Vendor: "Apple",
				Model:  "Apple M2 Pro GPU",
				Type:   "integrated",
			},
		},
		NPU: &internalmeta.NPUInfo{
			Detected:        true,
			Type:            "Apple Neural Engine",
			InferenceMethod: "cpu_model",
		},
	}

	output := formatSystemInfo(info)

	// Check OS section
	if !strings.Contains(output, "OS: darwin") {
		t.Error("missing OS field")
	}
	if !strings.Contains(output, "Platform: macOS 14.2") {
		t.Error("missing Platform field")
	}
	if !strings.Contains(output, "Architecture: arm64") {
		t.Error("missing Architecture field")
	}

	// Check CPU section
	if !strings.Contains(output, "CPU:") {
		t.Error("missing CPU section")
	}
	if !strings.Contains(output, "Model: Apple M2 Pro") {
		t.Error("missing CPU Model")
	}
	if !strings.Contains(output, "Cores: 10") {
		t.Error("missing CPU Cores")
	}

	// Check Memory section
	if !strings.Contains(output, "Memory:") {
		t.Error("missing Memory section")
	}
	if !strings.Contains(output, "Total: 16384 MB") {
		t.Error("missing Memory Total")
	}
	if !strings.Contains(output, "50.0%") {
		t.Error("missing Memory UsedPercent")
	}

	// Check Storage section
	if !strings.Contains(output, "Storage:") {
		t.Error("missing Storage section")
	}
	if !strings.Contains(output, "/: 505856 MB total") {
		t.Error("missing Storage mountpoint")
	}

	// Check GPU section
	if !strings.Contains(output, "GPU:") {
		t.Error("missing GPU section")
	}
	if !strings.Contains(output, "Apple M2 Pro GPU") {
		t.Error("missing GPU model")
	}

	// Check NPU section
	if !strings.Contains(output, "NPU:") {
		t.Error("missing NPU section")
	}
	if !strings.Contains(output, "Apple Neural Engine") {
		t.Error("missing NPU type")
	}
}

func TestFormatSystemInfo_NoGPU(t *testing.T) {
	info := internalmeta.SystemInfo{
		OS:           "linux",
		Platform:     "Ubuntu 22.04",
		Kernel:       "5.15.0",
		Architecture: "amd64",
		CPU: internalmeta.CPUInfo{
			Model:  "Intel Core i7",
			Vendor: "GenuineIntel",
			Cores:  8,
		},
		Memory:  internalmeta.MemoryInfo{TotalMB: 16384},
		Storage: []internalmeta.StorageInfo{},
		GPU:     []internalmeta.GPUInfo{},
		NPU:     nil,
	}

	output := formatSystemInfo(info)

	// Should have OS section
	if !strings.Contains(output, "OS: linux") {
		t.Error("missing OS field")
	}

	// Should NOT have GPU section (empty array)
	if strings.Contains(output, "GPU:") {
		t.Error("should not show GPU section when no GPUs detected")
	}

	// Should NOT have NPU section (nil)
	if strings.Contains(output, "NPU:") {
		t.Error("should not show NPU section when no NPU detected")
	}
}
