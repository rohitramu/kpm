package docker

import (
	dockerTypes "github.com/docker/docker/api/types"
)

func deleteContainer(docker dockerConnection, containerID string) error {
	var err error

	// Delete container
	var containerRemoveOpts = dockerTypes.ContainerRemoveOptions{
		Force:         true,
		RemoveLinks:   true,
		RemoveVolumes: true,
	}
	err = docker.Client.ContainerRemove(docker.Context, containerID, containerRemoveOpts)
	if err != nil {
		return err
	}

	return nil
}
