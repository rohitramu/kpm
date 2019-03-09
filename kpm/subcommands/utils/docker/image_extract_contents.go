package docker

import (
	"fmt"

	"../cmd"
	"../log"
)

// ExtractImageContents extracts files and directories from a Docker image, and copies them to the local KPM repository.
func ExtractImageContents(imageName string, destinationDir string) error {
	var err error

	var exe = "docker"
	const containerName = "kpm_container"

	// Create a container using the image
	{
		log.Info(fmt.Sprintf("Creating container from image: %s", imageName))

		var args = []string{"create", "--name", containerName, imageName}
		_, err = cmd.Exec(exe, args...)
		if err != nil {
			return fmt.Errorf("Failed to create container from image: %s\n%s", imageName, err)
		}
	}

	// Delete container after we're done
	defer func() {
		var args = []string{"rm", "--force", containerName}
		var deleteErr error
		_, deleteErr = cmd.Exec(exe, args...)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("%s\n%s", deleteErr, err)
			} else {
				err = deleteErr
			}

			err = fmt.Errorf("Failed to delete container: %s\n%s", containerName, err)
		}
	}()

	// Extract contents of container to a temporary directory
	{
		log.Info(fmt.Sprintf("Extracting contents from container: %s", containerName))

		var args = []string{"cp", fmt.Sprintf("%s:/%s", containerName, DockerfileRootDir), destinationDir}
		_, err = cmd.Exec(exe, args...)
		if err != nil {
			return fmt.Errorf("Failed to extract data from container: %s\n%s", containerName, err)
		}
	}

	return nil
}
