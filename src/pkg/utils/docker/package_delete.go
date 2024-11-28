package docker

import (
	"fmt"

	"github.com/rohitramu/kpm/src/pkg/utils/exec"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
)

// DeleteImage deletes a local Docker image.
func DeleteImage(imageName string) error {
	var err error

	log.Verbosef("Deleting image: %s", imageName)

	var exe = "docker"
	var args = []string{"image", "rm", "--force", imageName}
	_, err = exec.Exec(exe, args...)
	if err != nil {
		return fmt.Errorf("failed to delete image: %s\n%s", imageName, err)
	}

	return nil
}
