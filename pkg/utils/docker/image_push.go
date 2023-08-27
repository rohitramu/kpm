package docker

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/exec"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

// PushImage pushes a Docker image to a remote Docker registry.
func PushImage(imageName string) error {
	var err error

	log.Infof("Pushing image: %s", imageName)

	var exe = "docker"
	var args = []string{"push", imageName}
	_, err = exec.Exec(exe, args...)
	if err != nil {
		return fmt.Errorf("failed to push image: %s\n%s", imageName, err)
	}

	return nil
}
