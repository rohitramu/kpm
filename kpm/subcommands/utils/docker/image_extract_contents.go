package docker

import (
	"fmt"
	"strings"

	"../cmd"
	"../logger"
)

// ExtractImageContents extracts files and directories from a Docker image, and copies them to the local KPM repository.
func ExtractImageContents(imageName string, destinationDir string) error {
	var err error

	var exe = "docker"
	const containerName = "kpm_container"

	// Create a container using the image
	{
		logger.Default.Info.Println(fmt.Sprintf("Creating container from image: %s", imageName))

		var args = []string{"create", "--name", containerName, imageName}
		var output string
		output, err = cmd.Exec(exe, args...)
		logger.Default.Verbose.Println(fmt.Sprintf("%s %s\n%s", exe, strings.Join(args, " "), string(output)))
		if err != nil {
			return err
		}
	}

	// Delete container after we're done
	defer func() {
		var args = []string{"rm", "--force", containerName}
		var output string
		var deleteErr error
		output, deleteErr = cmd.Exec(exe, args...)
		logger.Default.Verbose.Println(fmt.Sprintf("%s %s\n%s", exe, strings.Join(args, " "), string(output)))
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to delete container:\n%s\n%s", deleteErr, err)
			} else {
				err = deleteErr
			}
		}
	}()

	// Extract contents of container to a temporary directory
	{
		logger.Default.Info.Println(fmt.Sprintf("Extracting contents from container: %s", containerName))

		var args = []string{"cp", fmt.Sprintf("%s:/%s", containerName, DockerfileRootDir), destinationDir}
		var output string
		output, err = cmd.Exec(exe, args...)
		logger.Default.Verbose.Println(fmt.Sprintf("%s %s\n%s", exe, strings.Join(args, " "), string(output)))
		if err != nil {
			return err
		}
	}

	return nil
}
