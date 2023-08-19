package main

import (
	"github.com/rohitramu/kpm/cli"
	"github.com/rohitramu/kpm/pkg/utils/config"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

// TODO: Make it possible to provide just the package name as an argument without the username prefix (unless it's ambiguous).
// TODO: Introduce the concept of "repositories".  The one with the highest priority will be the local repository.
//       Repositories can be added and removed by users.  Need to define different repository types (e.g. local filesystem, Docker, GitHub, etc.)
//       Need to handle authentication for private remote repositories (i.e. Docker, GitHub, etc.)

var VersionString = "0.0.0"

func main() {
	var err error

	err = cli.InitConfig(config.KpmConfig)
	if err != nil {
		log.Fatalf("Failed to read KPM configuration: %s", err)
	}

	if err = cli.RootCmd.Execute(); err != nil {
		log.Errorf("Failed to execute command: %s", err.Error())
	}
}
