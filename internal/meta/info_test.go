package meta

import (
	"runtime"
	"strings"
	"testing"
)

func TestCurrentBuildInfo(t *testing.T) {
	info := CurrentBuildInfo()

	// Name should always be set
	if info.Name != "ado" {
		t.Errorf("Name = %q, want %q", info.Name, "ado")
	}

	// Version, Commit, BuildTime are set via ldflags during build
	// In tests, they default to development values
	if info.Version == "" {
		t.Error("Version should not be empty")
	}
	if info.Commit == "" {
		t.Error("Commit should not be empty")
	}
	if info.BuildTime == "" {
		t.Error("BuildTime should not be empty")
	}

	// GoVersion should be populated
	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
	if !strings.HasPrefix(info.GoVersion, "go") {
		t.Errorf("GoVersion = %q, expected to start with 'go'", info.GoVersion)
	}

	// Platform should be in format "os/arch"
	if info.Platform == "" {
		t.Error("Platform should not be empty")
	}
	expectedPlatform := runtime.GOOS + "/" + runtime.GOARCH
	if info.Platform != expectedPlatform {
		t.Errorf("Platform = %q, want %q", info.Platform, expectedPlatform)
	}
}

func TestBuildInfo_Defaults(t *testing.T) {
	// Test that default build metadata values are set
	// These are used in development builds

	if Version == "" {
		t.Error("Version constant should have default value")
	}
	if Commit == "" {
		t.Error("Commit constant should have default value")
	}
	if BuildTime == "" {
		t.Error("BuildTime constant should have default value")
	}
}
