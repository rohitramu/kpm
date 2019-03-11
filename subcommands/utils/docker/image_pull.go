package docker

import (
	"fmt"

	"github.com/rohitramu/kpm/subcommands/utils/cmd"
	"github.com/rohitramu/kpm/subcommands/utils/log"
)

// PullImage pulls a Docker image from a remote Docker registry.
func PullImage(imageName string) error {
	var err error

	log.Verbose("Pulling image: %s", imageName)

	var exe = "docker"
	var args = []string{"pull", imageName}
	_, err = cmd.Exec(exe, args...)
	if err != nil {
		return fmt.Errorf("Failed to pull image: %s\n%s", imageName, err)
	}

	return nil
}
