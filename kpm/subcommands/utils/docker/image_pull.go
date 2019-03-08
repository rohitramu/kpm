package docker

import (
	"fmt"
	"strings"

	"../cmd"
	"../logger"
)

// PullImage pulls a Docker image from a remote Docker registry.
func PullImage(imageName string) error {
	logger.Default.Info.Println(fmt.Sprintf("Pulling image: %s", imageName))

	var exe = "docker"
	var args = []string{"pull", imageName}
	var output, err = cmd.Exec(exe, args...)
	logger.Default.Verbose.Println(fmt.Sprintf("%s %s\n%s", exe, strings.Join(args, " "), string(output)))
	if err != nil {
		return err
	}

	return nil
}
