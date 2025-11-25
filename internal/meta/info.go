package meta

import (
	"fmt"
	"runtime"
)

// Build metadata is set at compile time via -ldflags when available.
const Name = "ado"

var (
	Version   = "0.0.0-dev"
	Commit    = "none"
	BuildTime = "unknown"
)

type BuildInfo struct {
	Name      string `json:"name" yaml:"name"`
	Version   string `json:"version" yaml:"version"`
	Commit    string `json:"commit" yaml:"commit"`
	BuildTime string `json:"build_time" yaml:"build_time"`
	GoVersion string `json:"go_version" yaml:"go_version"`
	Platform  string `json:"platform" yaml:"platform"`
}

func CurrentBuildInfo() BuildInfo {
	return BuildInfo{
		Name:      Name,
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
