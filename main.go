package main

import (
	"github.com/rohitramu/kpm/src/cli"
	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/pkg/utils/log"

	"runtime/debug"
	"strconv"
)

// TODO: Make it possible to provide just the package name as an argument without the username prefix (unless it's ambiguous).
// TODO: Introduce the concept of "repositories".  The one with the highest priority will be the local repository.
//       Repositories can be added and removed by users.  Need to define different repository types (e.g. local filesystem, Docker, GitHub, etc.)
//       Need to handle authentication for private remote repositories (i.e. Docker, GitHub, etc.)

func main() {
	var err error

	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("Failed to get build info")
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			constants.GitCommitHash = setting.Value
		case "vcs.time":
			constants.BuildTimestampUTC = setting.Value
		case "vcs.modified":
			constants.IsSourceModified, _ = strconv.ParseBool(setting.Value)
		}
	}
	if constants.IsSourceModified {
		constants.VersionString += "-dirty"
		constants.GitCommitHash += " (modified)"
	}

	if err = cli.Execute(); err != nil {
		log.Errorf("Failed to execute command: %s", err.Error())
	}
}
