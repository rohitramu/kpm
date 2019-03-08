package docker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"

	"../files"
	"../logger"
)

// ExtractImageContents extracts files and directories from a Docker image, and copies them to the local KPM repository.
func ExtractImageContents(imageName string, destinationDir string) error {
	var err error

	logger.Default.Info.Println(fmt.Sprintf("Extracting image contents to: %s", destinationDir))

	// Make sure imageDir is a Unix absolute file path
	const imageDir = DockerTarFileRootDir

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
		WorkingDir:      "/",
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
	logger.Default.Verbose.Println(fmt.Sprintf("Created new container: %s\nID: %s\nWarnings: %s", imageName, createResponse.ID, createResponse.Warnings))

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
	defer func() {
		var closeErr = tarData.Close()
		if closeErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to close tar data stream:\n%s\n%s", closeErr, err)
			} else {
				err = closeErr
			}
		}
	}()

	// We want to extract files from tar data to a temporary location first in case of a failure, so clear the temporary directory
	var imageNameWithoutColon = strings.Replace(imageName, ":", "-", -1)
	var tempDir = filepath.Join(os.TempDir(), ".kpm", imageNameWithoutColon)
	logger.Default.Verbose.Println(fmt.Sprintf("Extracting to temporary directory: %s", tempDir))
	err = os.RemoveAll(tempDir)
	if err != nil {
		logger.Default.Error.Panicln(err)
	}
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to make temporary directory for extracting files: %s", err)
	}

	// Extract tar file to temporary directory
	err = extractTar(tarData, imageDir, tempDir)
	if err != nil {
		return fmt.Errorf("Failed to extract tar file from Docker image: %s", err)
	}

	// Delete destination directory and recreate it in case already exists
	logger.Default.Verbose.Println(fmt.Sprintf("Clearing destination directory: %s", destinationDir))
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

	// // Make sure to delete the temporary directory once we're done if we were successful
	// defer func() {
	// 	logger.Default.Verbose.Println(fmt.Sprintf("Removing temporary directory: %s", tempDir))
	// 	var dirExists = files.DirExists(tempDir, "temporary") == nil
	// 	if dirExists {
	// 		var deleteErr = os.RemoveAll(tempDir)
	// 		if deleteErr != nil {
	// 			if err != nil {
	// 				err = fmt.Errorf("Failed to clean up temporary directory:\n%s\n%s", deleteErr, err)
	// 			} else {
	// 				err = deleteErr
	// 			}
	// 		}
	// 	}
	// }()

	// Handle copy error
	if err != nil {
		return err
	}

	return err
}
