package docker

import (
	"fmt"

	"github.com/rohitramu/kpm/subcommands/utils/cmd"
	"github.com/rohitramu/kpm/subcommands/utils/log"
)

// DeleteImage deletes a local Docker image.
func DeleteImage(imageName string) error {
	var err error

	log.Verbose("Deleting image: %s", imageName)

	var exe = "docker"
	var args = []string{"image", "rm", "--force", imageName}
	_, err = cmd.Exec(exe, args...)
	if err != nil {
		return fmt.Errorf("Failed to delete image: %s\n%s", imageName, err)
	}

	return nil
}
