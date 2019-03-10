package docker

import (
	"fmt"

	"../cmd"
	"../log"
)

// PushImage pushes a Docker image to a remote Docker registry.
func PushImage(imageName string) error {
	var err error

	log.Info("Pushing image: %s", imageName)

	var exe = "docker"
	var args = []string{"push", imageName}
	_, err = cmd.Exec(exe, args...)
	if err != nil {
		return fmt.Errorf("Failed to push image: %s\n%s", imageName, err)
	}

	return nil
}
