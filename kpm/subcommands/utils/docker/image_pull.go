package docker

import (
	"fmt"
	"io"
	"strings"

	"../logger"
	"./credentials"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

// PullImage pulls a Docker image from a remote Docker registry.
func PullImage(dockerRegistryURL string, imageName string) error {
	var err error

	// Get Docker client
	var docker dockerConnection
	docker, err = getClient()
	if err != nil {
		return err
	}

	// Get Docker credentials
	var authString string
	authString, err = credentials.GetCredentialsFromConfig(dockerRegistryURL)
	if err != nil {
		return err
	}

	// Create the pull options
	var pullOpts = dockerTypes.ImagePullOptions{
		All:          true,
		RegistryAuth: authString,
	}

	// Pull image
	var pullResponse io.ReadCloser
	pullResponse, err = docker.Client.ImagePull(docker.Context, imageName, pullOpts)
	if err != nil {
		return fmt.Errorf("Failed to pull image \"%s\": %s", imageName, err)
	}
	defer pullResponse.Close()

	// Print output
	var stringStream = &strings.Builder{}
	termFd, isTerm := term.GetFdInfo(stringStream)
	err = jsonmessage.DisplayJSONMessagesStream(pullResponse, stringStream, termFd, isTerm, nil)
	if err != nil {
		return err
	}
	logger.Default.Verbose.Println(fmt.Sprintf("Docker build:\n%s", stringStream.String()))

	return nil
}
