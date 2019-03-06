package docker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"

	"../files"
	"../logger"
)

// CopyImageContents extracts files and directories from a Docker image, and copies them to the destination directory.
func CopyImageContents(imageName string, imageDir string, destinationDir string) error {
	var err error

	// Make sure imageDir is a Unix file path
	imageDir = filepath.ToSlash(imageDir)

	// If imageDir is a relative path, make it an absolute path
	if !filepath.IsAbs(imageDir) {
		imageDir = "/" + imageDir
	}

	// Get Docker client
	var docker dockerConnection
	docker, err = getClient()
	if err != nil {
		return err
	}

	// Create container
	logger.Default.Verbose.Println(fmt.Sprintf("Creating container for image: %s", imageName))
	var containerConfig = &container.Config{
		Image:           imageName,
		WorkingDir:      imageDir,
		Cmd:             []string{""},
		NetworkDisabled: true,
	}
	var hostConfig = &container.HostConfig{}
	var networkingConfig = &network.NetworkingConfig{}
	var createResponse container.ContainerCreateCreatedBody
	createResponse, err = docker.Client.ContainerCreate(docker.Context, containerConfig, hostConfig, networkingConfig, "")
	if err != nil {
		return fmt.Errorf("Failed to create Docker container:\n%s", err)
	}

	// Delete the container after we're done
	defer func() {
		logger.Default.Verbose.Println(fmt.Sprintf("Deleting Docker container: %s", createResponse.ID))
		// First check to see if the container exists
		var deleteErr = deleteContainer(docker, createResponse.ID)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to clean up container:\n%s\n%s", deleteErr, err)
			} else {
				err = deleteErr
			}
		}
	}()

	// Get tar data from container
	logger.Default.Verbose.Println(fmt.Sprintf("Reading contents of container: %s", createResponse.ID))
	var tarData io.ReadCloser
	tarData, _, err = docker.Client.CopyFromContainer(docker.Context, createResponse.ID, imageDir)
	if err != nil {
		logger.Default.Error.Println(fmt.Sprintf("Failed to copy package from container: %s", createResponse.ID))
		return err
	}
	defer tarData.Close()

	// Extract files from tar data to a temporary location first, in case of a failure
	var tempDir = filepath.Join(os.TempDir(), ".kpm", filepath.Base(destinationDir))
	logger.Default.Verbose.Println(fmt.Sprintf("Extracting contents of container to: %s", tempDir))
	err = os.RemoveAll(tempDir)
	if err != nil {
		logger.Default.Error.Panicln(err)
	}
	err = extractTar(tarData, tempDir)
	if err != nil {
		return fmt.Errorf("Failed to extract tar file from Docker image:\n%s", err)
	}

	// Delete destination directory and recreate it in case already exists
	logger.Default.Verbose.Println(fmt.Sprintf("Re-creating destination directory: %s", destinationDir))
	err = os.RemoveAll(destinationDir)
	if err != nil {
		logger.Default.Error.Panicln(err)
	}
	err = os.MkdirAll(destinationDir, os.ModePerm)
	if err != nil {
		logger.Default.Error.Panicln(err)
	}

	// Copy files from the temporary location to the destination directory
	logger.Default.Verbose.Println(fmt.Sprintf("Copying package to local repository: %s", destinationDir))
	err = files.CopyDir(tempDir, destinationDir)

	// Make sure to delete the temporary directory once we're done
	defer func() {
		logger.Default.Verbose.Println(fmt.Sprintf("Removing temporary directory: %s", tempDir))
		var deleteErr = os.RemoveAll(tempDir)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to clean up temporary directory:\n%s\n%s", deleteErr, err)
			} else {
				err = deleteErr
			}
		}

		// We should never fail here, so panic if we did
		if err != nil {
			logger.Default.Error.Panicln(err)
		}
	}()

	// Handle copy error
	if err != nil {
		return err
	}

	return nil
}
