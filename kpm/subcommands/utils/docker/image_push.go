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

// PushImage pushes a Docker image to a remote Docker registry.
func PushImage(dockerRegistryURL string, imageName string) error {
	var err error

	logger.Default.Info.Println(fmt.Sprintf("Pushing image: %s", imageName))

	// Get Docker client
	var docker dockerConnection
	docker, err = getClient()
	if err != nil {
		return err
	}

	// Get Docker config
	var config *credentials.DockerConfig
	config, err = credentials.GetDockerConfig()
	if err != nil {
		return err
	}

	// Get Docker credentials
	var authString string
	authString, err = config.GetCredentials(dockerRegistryURL)
	if err != nil {
		return err
	}

	// Create the push options
	var pushOpts = dockerTypes.ImagePushOptions{
		All:          true,
		RegistryAuth: authString,
	}

	// Push image
	var pushResponse io.ReadCloser
	pushResponse, err = docker.Client.ImagePush(docker.Context, imageName, pushOpts)
	if err != nil {
		return err
	}
	defer pushResponse.Close()

	// Print output
	var stringStream = &strings.Builder{}
	termFd, isTerm := term.GetFdInfo(stringStream)
	err = jsonmessage.DisplayJSONMessagesStream(pushResponse, stringStream, termFd, isTerm, nil)
	if err != nil {
		return err
	}
	logger.Default.Verbose.Println(fmt.Sprintf("Docker push:\n%s", stringStream.String()))

	return nil
}
