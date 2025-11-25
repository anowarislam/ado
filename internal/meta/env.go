package meta

import (
	"os"

	"github.com/anowarislam/ado/internal/config"
)

type EnvInfo struct {
	ConfigPath    string            `json:"config_path" yaml:"config_path"`
	ConfigSources []string          `json:"config_sources" yaml:"config_sources"`
	HomeDir       string            `json:"home_dir" yaml:"home_dir"`
	CacheDir      string            `json:"cache_dir" yaml:"cache_dir"`
	Env           map[string]string `json:"env" yaml:"env"`
}

func CollectEnvInfo(explicitConfig string) EnvInfo {
	homeDir, _ := os.UserHomeDir()
	cacheDir, _ := os.UserCacheDir()

	configPath := explicitConfig
	if configPath == "" {
		if envPath, ok := os.LookupEnv("ADO_CONFIG"); ok {
			configPath = envPath
		}
	}

	resolved, sources := config.ResolveConfigPath(configPath, homeDir)

	envVars := map[string]string{}
	for _, key := range []string{"ADO_CONFIG", "ADO_LOG_LEVEL"} {
		if value, ok := os.LookupEnv(key); ok {
			envVars[key] = value
		}
	}

	return EnvInfo{
		ConfigPath:    resolved,
		ConfigSources: sources,
		HomeDir:       homeDir,
		CacheDir:      cacheDir,
		Env:           envVars,
	}
}
