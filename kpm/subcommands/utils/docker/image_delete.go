package docker

import (
	"fmt"
	"strings"

	"../cmd"
	"../logger"
)

// DeleteImage deletes a local Docker image.
func DeleteImage(imageName string) error {
	logger.Default.Info.Println(fmt.Sprintf("Deleting image: %s", imageName))

	var exe = "docker"
	var args = []string{"image", "rm", "--force", imageName}
	var output, err = cmd.Exec(exe, args...)
	logger.Default.Verbose.Println(fmt.Sprintf("%s %s\n%s", exe, strings.Join(args, " "), string(output)))
	if err != nil {
		return err
	}

	return nil
}
