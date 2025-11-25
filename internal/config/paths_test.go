package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefaultSearchPaths_XDGPreferred(t *testing.T) {
	home := t.TempDir()
	xdg := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_CONFIG_HOME", xdg)

	got := DefaultSearchPaths(home)
	want := []string{
		filepath.Join(xdg, "ado", "config.yaml"),
		filepath.Join(home, ".ado", "config.yaml"),
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DefaultSearchPaths mismatch\n  got:  %#v\n  want: %#v", got, want)
	}
}

func TestResolveConfigPath_FindsXDGConfig(t *testing.T) {
	home := t.TempDir()
	xdg := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_CONFIG_HOME", xdg)

	xdgConfig := filepath.Join(xdg, "ado", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(xdgConfig), 0o755); err != nil {
		t.Fatalf("mkdir xdg config dir: %v", err)
	}
	if err := os.WriteFile(xdgConfig, []byte("content"), 0o644); err != nil {
		t.Fatalf("write xdg config: %v", err)
	}

	gotPath, gotSources := ResolveConfigPath("", home)

	wantSources := []string{
		filepath.Join(xdg, "ado", "config.yaml"),
		filepath.Join(home, ".ado", "config.yaml"),
	}

	if gotPath != xdgConfig {
		t.Fatalf("ResolveConfigPath path mismatch: got %q want %q", gotPath, xdgConfig)
	}
	if !reflect.DeepEqual(gotSources, wantSources) {
		t.Fatalf("ResolveConfigPath sources mismatch\n  got:  %#v\n  want: %#v", gotSources, wantSources)
	}
}

func TestResolveConfigPath_NoConfigFound(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")

	gotPath, gotSources := ResolveConfigPath("", home)

	wantSources := []string{
		filepath.Join(home, ".config", "ado", "config.yaml"),
		filepath.Join(home, ".ado", "config.yaml"),
	}

	if gotPath != "" {
		t.Fatalf("expected no config path, got %q", gotPath)
	}
	if !reflect.DeepEqual(gotSources, wantSources) {
		t.Fatalf("ResolveConfigPath sources mismatch\n  got:  %#v\n  want: %#v", gotSources, wantSources)
	}
}

func TestResolveConfigPath_ExplicitPathWins(t *testing.T) {
	home := t.TempDir()
	xdg := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_CONFIG_HOME", xdg)

	explicit := filepath.Join(t.TempDir(), "custom-config.yaml")

	gotPath, gotSources := ResolveConfigPath(explicit, home)

	wantSources := []string{
		explicit,
		filepath.Join(xdg, "ado", "config.yaml"),
		filepath.Join(home, ".ado", "config.yaml"),
	}

	if gotPath != explicit {
		t.Fatalf("ResolveConfigPath path mismatch: got %q want %q", gotPath, explicit)
	}
	if !reflect.DeepEqual(gotSources, wantSources) {
		t.Fatalf("ResolveConfigPath sources mismatch\n  got:  %#v\n  want: %#v", gotSources, wantSources)
	}
}
