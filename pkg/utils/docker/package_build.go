package docker

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/exec"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

// BuildImage builds a new docker image by making a call to the Docker daemon.
func BuildImage(imageName string, dockerfilePath string, dirToCopy string) error {
	log.Infof("Building image: %s", imageName)

	var exe = "docker"
	var args = []string{"build", "--force-rm", "--file", dockerfilePath, "--tag", imageName, dirToCopy}
	var _, err = exec.Exec(exe, args...)
	if err != nil {
		return fmt.Errorf("failed to build image: %s\n%s", imageName, err)
	}

	return nil
}
