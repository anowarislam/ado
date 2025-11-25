package config

import (
	"os"
	"path/filepath"
)

// DefaultSearchPaths returns the default config lookup order, excluding any explicit flag value.
func DefaultSearchPaths(homeDir string) []string {
	var paths []string

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "ado", "config.yaml"))
	} else if homeDir != "" {
		paths = append(paths, filepath.Join(homeDir, ".config", "ado", "config.yaml"))
	}

	if homeDir != "" {
		paths = append(paths, filepath.Join(homeDir, ".ado", "config.yaml"))
	}

	return paths
}

// ResolveConfigPath returns the resolved config path (if found) and the list of sources checked.
func ResolveConfigPath(explicitPath, homeDir string) (string, []string) {
	if explicitPath != "" {
		return explicitPath, append([]string{explicitPath}, DefaultSearchPaths(homeDir)...)
	}

	sources := DefaultSearchPaths(homeDir)
	for _, candidate := range sources {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, sources
		}
	}

	return "", sources
}
