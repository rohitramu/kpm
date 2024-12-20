package docker

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rohitramu/kpm/src/pkg/utils/exec"
	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
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
		log.Debugf("Creating container \"%s\" from image: %s", containerName, imageName)

		var args = []string{"create", "--name", containerName, imageName}
		_, err = exec.Exec(exe, args...)
		if err != nil {
			return fmt.Errorf("failed to create container from image: %s\n%s", imageName, err)
		}
	}

	// Delete container after we're done
	defer func() {
		log.Debugf("Deleting container: %s", containerName)

		var args = []string{"rm", "--force", containerName}
		var deleteErr error
		_, deleteErr = exec.Exec(exe, args...)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("%s\n%s", deleteErr, err)
			} else {
				err = deleteErr
			}

			err = fmt.Errorf("failed to delete container: %s\n%s", containerName, err)
		}
	}()

	// Extract contents of container to a temporary directory
	var imageNameWithoutColon = strings.Replace(imageName, ":", "-", -1)
	var tempDir string
	tempDir, err = files.GetTempDir()
	if err != nil {
		return err
	}
	tempDir = filepath.Join(tempDir, ".kpm", imageNameWithoutColon)
	{
		log.Debugf("Extracting contents from container \"%s\" to temporary directory: %s", containerName, tempDir)

		// Remove temporary directory to clear it
		err = files.DeleteDirIfExists(tempDir, "temporary", true)
		if err != nil {
			log.Panicf("Failed to remove directory: %s", err)
		}

		// Recreate temporary directory
		err = files.CreateDir(tempDir, "temporary", true)
		if err != nil {
			return fmt.Errorf("failed to make temporary directory for extracting files: %s\n%s", tempDir, err)
		}

		// Extract data to temporary directory
		var args = []string{"cp", fmt.Sprintf("%s:/%s/.", containerName, DockerfileRootDir), tempDir}
		_, err = exec.Exec(exe, args...)
		if err != nil {
			return fmt.Errorf("failed to extract data from container: %s\n%s", containerName, err)
		}
	}

	// Copy data to destination directory
	{
		log.Debugf("Copying contents of container \"%s\" to destination directory: %s", containerName, destinationDir)

		// Remove destination directory to clear it
		err = files.DeleteDirIfExists(destinationDir, "template package", true)
		if err != nil {
			log.Panicf("Failed to remove directory: %s", err)
		}

		// Recreate destination directory
		err = files.CreateDir(destinationDir, "template package", true)
		if err != nil {
			return fmt.Errorf("failed to make destination directory for saving data: %s\n%s", destinationDir, err)
		}

		// Copy data to destination directory
		err = files.CopyDir(tempDir, destinationDir)
		if err != nil {
			return fmt.Errorf("failed to copy data from temporary directory to destination directory: %s -> %s\n%s", tempDir, destinationDir, err)
		}
	}

	// Delete temporary directory
	{
		log.Debugf("Deleting temporary directory: %s", tempDir)
		err = files.DeleteDirIfExists(tempDir, "temporary", true)
		if err != nil {
			return fmt.Errorf("failed to delete temporary directory: %s", err)
		}
	}

	return nil
}
