package docker

import (
	"fmt"

	"github.com/rohitramu/kpm/subcommands/utils/cmd"
	"github.com/rohitramu/kpm/subcommands/utils/log"
)

// BuildImage builds a new docker image by making a call to the Docker daemon.
func BuildImage(imageName string, dockerfilePath string, dirToCopy string) error {
	log.Info("Building image: %s", imageName)

	var exe = "docker"
	var args = []string{"build", "--force-rm", "--file", dockerfilePath, "--tag", imageName, dirToCopy}
	var _, err = cmd.Exec(exe, args...)
	if err != nil {
		return fmt.Errorf("Failed to build image: %s\n%s", imageName, err)
	}

	return nil
}
