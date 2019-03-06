package docker

import (
	"bytes"
	"fmt"
	"strings"

	"../logger"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

// BuildImage builds a new docker image by making a call to the Docker daemon.
func BuildImage(imageName string, dockerfile string, dirToCopy string) error {
	var err error

	// Get Docker client
	var docker dockerConnection
	docker, err = getClient()
	if err != nil {
		return err
	}

	// Create build options
	var dockerfileName = "Dockerfile"
	var buildOptions = dockerTypes.ImageBuildOptions{
		Tags:           []string{imageName},
		Dockerfile:     dockerfileName,
		ForceRemove:    true,
		Remove:         true,
		PullParent:     true,
		SuppressOutput: false,
	}

	// Create in-memory tar file to use as the body in the request to the docker daemon
	var buildRequestBytes *bytes.Buffer
	buildRequestBytes, err = createTar(dockerfileName, dockerfile, dirToCopy)
	if err != nil {
		return err
	}

	// Send request to the docker daemon to build the image
	var response dockerTypes.ImageBuildResponse
	response, err = docker.Client.ImageBuild(docker.Context, buildRequestBytes, buildOptions)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Print output
	logger.Default.Verbose.Println(fmt.Sprintf("Docker daemon's reported OS: %s", response.OSType))
	var stringStream = &strings.Builder{}
	termFd, isTerm := term.GetFdInfo(stringStream)
	err = jsonmessage.DisplayJSONMessagesStream(response.Body, stringStream, termFd, isTerm, nil)
	if err != nil {
		return err
	}
	logger.Default.Verbose.Println(fmt.Sprintf("Docker build:\n%s", stringStream.String()))

	return nil
}
