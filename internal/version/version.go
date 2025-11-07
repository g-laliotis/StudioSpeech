// Package version provides build information and version details for StudioSpeech
package version

import (
	"fmt"
	"runtime"
)

// Build information set by ldflags during compilation
var (
	// Version is the semantic version of the application
	Version = "dev"
	
	// BuildTime is when the binary was built
	BuildTime = "unknown"
	
	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// Info contains version and build information
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get returns version information
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetBuildInfo returns formatted build information string
func GetBuildInfo() string {
	info := Get()
	return fmt.Sprintf("StudioSpeech %s\nCommit: %s\nBuilt: %s\nGo: %s\nPlatform: %s",
		info.Version,
		info.GitCommit,
		info.BuildTime,
		info.GoVersion,
		info.Platform,
	)
}

// GetVersion returns just the version string
func GetVersion() string {
	return Version
}

// IsDevBuild returns true if this is a development build
func IsDevBuild() bool {
	return Version == "dev"
}