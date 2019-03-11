package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rohitramu/kpm/subcommands/utils/cmd"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
)

// ExtractImageContents extracts files and directories from a Docker image, and copies them to the local KPM repository.
func ExtractImageContents(imageName string, destinationDir string) error {
	var err error

	const (
		exe           = "docker"
		containerName = "kpm_container"
	)

	// Create a container using the image
	{
		log.Verbose("Creating container \"%s\" from image: %s", containerName, imageName)

		var args = []string{"create", "--name", containerName, imageName}
		_, err = cmd.Exec(exe, args...)
		if err != nil {
			return fmt.Errorf("Failed to create container from image: %s\n%s", imageName, err)
		}
	}

	// Delete container after we're done
	defer func() {
		log.Verbose("Deleting container: %s", containerName)

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
	var imageNameWithoutColon = strings.Replace(imageName, ":", "-", -1)
	var tempDir = filepath.Join(os.TempDir(), ".kpm", imageNameWithoutColon)
	{
		log.Verbose("Extracting contents from container \"%s\" to temporary directory: %s", containerName, tempDir)

		// Remove temporary directory to clear it
		err = os.RemoveAll(tempDir)
		if err != nil {
			log.Panic("Failed to remove directory: %s", err)
		}

		// Recreate temporary directory
		err = os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to make temporary directory for extracting files: %s\n%s", tempDir, err)
		}

		// Extract data to temporary directory
		var args = []string{"cp", fmt.Sprintf("%s:/%s/.", containerName, DockerfileRootDir), tempDir}
		_, err = cmd.Exec(exe, args...)
		if err != nil {
			return fmt.Errorf("Failed to extract data from container: %s\n%s", containerName, err)
		}
	}

	// Copy data to destination directory
	{
		log.Verbose("Copying contents of container \"%s\" to destination directory: %s", containerName, destinationDir)

		// Remove destination directory to clear it
		err = os.RemoveAll(destinationDir)
		if err != nil {
			log.Panic("Failed to remove directory: %s", err)
		}

		// Recreate destination directory
		err = os.MkdirAll(destinationDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Failed to make destination directory for saving data: %s\n%s", destinationDir, err)
		}

		// Copy data to destination directory
		err = files.CopyDir(tempDir, destinationDir)
		if err != nil {
			return fmt.Errorf("Failed to copy data from temporary directory to destination directory: %s -> %s\n%s", tempDir, destinationDir, err)
		}
	}

	// Delete temporary directory
	{
		log.Verbose("Deleting temporary directory: %s", tempDir)
		err = os.RemoveAll(tempDir)
		if err != nil {
			return fmt.Errorf("Failed to delete temporary directory: %s", err)
		}
	}

	return nil
}
