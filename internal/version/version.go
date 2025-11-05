package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the semantic version of the application
	Version = "dev"
	// BuildTime is the time the binary was built
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// GitBranch is the git branch
	GitBranch = "unknown"
)

// Info represents version information
type Info struct {
	Version   string
	BuildTime string
	GitCommit string
	GitBranch string
	GoVersion string
	Platform  string
}

// Get returns the version information
func Get() Info {
	return Info{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GitBranch: GitBranch,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf(
		"Unlinked %s\n"+
			"  Build Time: %s\n"+
			"  Git Commit: %s\n"+
			"  Git Branch: %s\n"+
			"  Go Version: %s\n"+
			"  Platform:   %s",
		i.Version,
		i.BuildTime,
		i.GitCommit,
		i.GitBranch,
		i.GoVersion,
		i.Platform,
	)
}

// Short returns a short version string
func (i Info) Short() string {
	if i.GitCommit != "unknown" && len(i.GitCommit) > 7 {
		return fmt.Sprintf("%s (%s)", i.Version, i.GitCommit[:7])
	}
	return i.Version
}
