package version

import "fmt"

const (
	Major = 1
	Minor = 0
	Patch = 0
)

var (
	Version   = fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
	GitCommit = "unknown"
	BuildDate = "unknown"
	GoVersion = "unknown"
)

func Full() string {
	return fmt.Sprintf("p2p-chat version %s\nGit commit: %s\nBuild date: %s\nGo version: %s",
		Version, GitCommit, BuildDate, GoVersion)
}

func Short() string {
	return Version
}

func BuildInfo() (string, string, string) {
	return GitCommit, BuildDate, GoVersion
}