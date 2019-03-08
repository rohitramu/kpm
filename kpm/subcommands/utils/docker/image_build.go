package docker

import (
	"fmt"
	"strings"

	"../cmd"
	"../logger"
)

// BuildImage builds a new docker image by making a call to the Docker daemon.
func BuildImage(imageName string, dockerfilePath string, dirToCopy string) error {
	logger.Default.Info.Println(fmt.Sprintf("Building image: %s", imageName))

	var exe = "docker"
	var args = []string{"build", "--force-rm", "--file", dockerfilePath, "--tag", imageName, dirToCopy}
	var output, err = cmd.Exec(exe, args...)
	logger.Default.Verbose.Println(fmt.Sprintf("%s %s:\n%s", exe, strings.Join(args, " "), string(output)))
	if err != nil {
		return err
	}

	return nil
}
