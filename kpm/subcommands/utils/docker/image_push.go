package docker

import (
	"fmt"
	"strings"

	"../cmd"
	"../logger"
)

// PushImage pushes a Docker image to a remote Docker registry.
func PushImage(imageName string) error {
	logger.Default.Info.Println(fmt.Sprintf("Pushing image: %s", imageName))

	var exe = "docker"
	var args = []string{"push", imageName}
	var output, err = cmd.Exec(exe, args...)
	logger.Default.Verbose.Println(fmt.Sprintf("%s %s\n%s", exe, strings.Join(args, " "), string(output)))
	if err != nil {
		return err
	}

	return nil
}
