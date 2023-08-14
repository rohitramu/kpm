package main

import (
	"github.com/rohitramu/kpm/cmd"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

// TODO: Make it possible to provide just the package name as an argument without the username prefix (unless it's ambiguous).
// TODO: Introduce the concept of "repositories".  The one with the highest priority will be the local repository.
//       Repositories can be added and removed by users.  Need to define different repository types (e.g. local filesystem, Docker, GitHub, etc.)
//       Need to handle authentication for private remote repositories (i.e. Docker, GitHub, etc.)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Errorf("Failed to execute command: %s", err.Error())
	}
}
