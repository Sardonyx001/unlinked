package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GitBranch = "unknown"
)

type Info struct {
	Version   string
	BuildTime string
	GitCommit string
	GitBranch string
	GoVersion string
	Platform  string
}

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

func (i Info) Short() string {
	if i.GitCommit != "unknown" && len(i.GitCommit) > 7 {
		return fmt.Sprintf("%s (%s)", i.Version, i.GitCommit[:7])
	}
	return i.Version
}
