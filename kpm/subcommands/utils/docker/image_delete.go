package docker

import (
	"fmt"

	"../logger"
	dockerTypes "github.com/docker/docker/api/types"
)

// DeleteImage deletes a local Docker image.
func DeleteImage(imageName string) error {
	var err error

	logger.Default.Info.Println(fmt.Sprintf("Deleting image: %s", imageName))

	// Get Docker client
	var docker dockerConnection
	docker, err = getClient()
	if err != nil {
		return err
	}

	// Set options
	var removeOpts = dockerTypes.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}

	// Delete the image
	var removeResponse []dockerTypes.ImageDeleteResponseItem
	removeResponse, err = docker.Client.ImageRemove(docker.Context, imageName, removeOpts)
	if err != nil {
		return err
	}

	// Print output
	for _, r := range removeResponse {
		if r.Untagged != "" {
			logger.Default.Verbose.Println(fmt.Sprintf("Untagged: %s", r.Untagged))
		}

		if r.Deleted != "" {
			logger.Default.Verbose.Println(fmt.Sprintf("Deleted:  %s", r.Deleted))
		}
	}

	return nil
}
