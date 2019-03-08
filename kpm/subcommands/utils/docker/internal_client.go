package docker

import (
	"context"

	"github.com/docker/docker/client"
)

type dockerConnection struct {
	Client  client.Client
	Context context.Context
}

var internalDockerConnection *dockerConnection

// getClient creates a new Docker request context and client.
func getClient() (dockerConnection, error) {
	var err error

	// Check if we already have a connection available
	if internalDockerConnection != nil {
		return *internalDockerConnection, nil
	}

	// Get a context object for making requests to the Docker daemon
	var result = dockerConnection{
		Context: context.Background(),
	}

	// Get Docker client
	var docker *client.Client
	docker, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return result, err
	}

	// Negotiate Docker daemon version
	docker.NegotiateAPIVersion(result.Context)

	// Set the Docker client
	result.Client = *docker
	internalDockerConnection = &result

	return result, nil
}
