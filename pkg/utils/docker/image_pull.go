package docker

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/exec"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

// PullImage pulls a Docker image from a remote Docker registry.
func PullImage(imageName string) error {
	var err error

	log.Verbosef("Pulling image: %s", imageName)

	var exe = "docker"
	var args = []string{"pull", imageName}
	_, err = exec.Exec(exe, args...)
	if err != nil {
		return fmt.Errorf("failed to pull image: %s\n%s", imageName, err)
	}

	return nil
}
