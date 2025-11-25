package meta

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCollectEnvInfo_ExplicitConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	xdg := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_CONFIG_HOME", xdg)

	explicit := filepath.Join(t.TempDir(), "custom-config.yaml")
	if err := os.WriteFile(explicit, []byte{}, 0o644); err != nil {
		t.Fatalf("write explicit config: %v", err)
	}

	t.Setenv("ADO_CONFIG", "/env/config.yaml")
	t.Setenv("ADO_LOG_LEVEL", "debug")

	info := CollectEnvInfo(explicit)

	wantSources := []string{
		explicit,
		filepath.Join(xdg, "ado", "config.yaml"),
		filepath.Join(home, ".ado", "config.yaml"),
	}

	if info.ConfigPath != explicit {
		t.Fatalf("ConfigPath mismatch: got %q want %q", info.ConfigPath, explicit)
	}
	if !reflect.DeepEqual(info.ConfigSources, wantSources) {
		t.Fatalf("ConfigSources mismatch\n  got:  %#v\n  want: %#v", info.ConfigSources, wantSources)
	}
	if info.HomeDir != home {
		t.Fatalf("HomeDir mismatch: got %q want %q", info.HomeDir, home)
	}

	expectedCacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatalf("resolve cache dir: %v", err)
	}
	if info.CacheDir != expectedCacheDir {
		t.Fatalf("CacheDir mismatch: got %q want %q", info.CacheDir, expectedCacheDir)
	}

	wantEnv := map[string]string{
		"ADO_CONFIG":    "/env/config.yaml",
		"ADO_LOG_LEVEL": "debug",
	}
	if !reflect.DeepEqual(info.Env, wantEnv) {
		t.Fatalf("Env mismatch\n  got:  %#v\n  want: %#v", info.Env, wantEnv)
	}
}

func TestCollectEnvInfo_DefaultResolution(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	configPath := filepath.Join(home, ".config", "ado", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	if err := os.WriteFile(configPath, []byte{}, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	info := CollectEnvInfo("")

	wantSources := []string{
		filepath.Join(home, ".config", "ado", "config.yaml"),
		filepath.Join(home, ".ado", "config.yaml"),
	}

	if info.ConfigPath != configPath {
		t.Fatalf("ConfigPath mismatch: got %q want %q", info.ConfigPath, configPath)
	}
	if !reflect.DeepEqual(info.ConfigSources, wantSources) {
		t.Fatalf("ConfigSources mismatch\n  got:  %#v\n  want: %#v", info.ConfigSources, wantSources)
	}
	if info.HomeDir != home {
		t.Fatalf("HomeDir mismatch: got %q want %q", info.HomeDir, home)
	}
	expectedCacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatalf("resolve cache dir: %v", err)
	}
	if info.CacheDir != expectedCacheDir {
		t.Fatalf("CacheDir mismatch: got %q want %q", info.CacheDir, expectedCacheDir)
	}
	if len(info.Env) != 0 {
		t.Fatalf("expected no env entries, got %#v", info.Env)
	}
}

func TestCollectEnvInfo_AdoConfigEnvPreferred(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	xdg := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_CONFIG_HOME", xdg)

	envConfig := filepath.Join(t.TempDir(), "env-config.yaml")
	t.Setenv("ADO_CONFIG", envConfig)

	info := CollectEnvInfo("")

	wantSources := []string{
		envConfig,
		filepath.Join(xdg, "ado", "config.yaml"),
		filepath.Join(home, ".ado", "config.yaml"),
	}

	if info.ConfigPath != envConfig {
		t.Fatalf("ConfigPath mismatch: got %q want %q", info.ConfigPath, envConfig)
	}
	if !reflect.DeepEqual(info.ConfigSources, wantSources) {
		t.Fatalf("ConfigSources mismatch\n  got:  %#v\n  want: %#v", info.ConfigSources, wantSources)
	}
	if val, ok := info.Env["ADO_CONFIG"]; !ok || val != envConfig {
		t.Fatalf("expected ADO_CONFIG to be captured, got %#v", info.Env)
	}
}
