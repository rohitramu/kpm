package docker

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/mholt/archiver"

	"../logger"
	"../types"
)

// DefaultDockerRegistryURL is the default registry URL (Docker Hub)
const DefaultDockerRegistryURL = "docker.io"

// GetClient creates a new Docker client.
func GetClient(dockerRegistryURL string) (context.Context, *client.Client, error) {
	var err error

	// Get a context object for making requests to the Docker daemon
	var dockerContext = context.Background()

	// Get docker client
	var docker *client.Client
	docker, err = client.NewEnvClient()
	if err != nil {
		return nil, nil, err
	}

	// Negotiate Docker daemon version
	docker.NegotiateAPIVersion(dockerContext)

	// // Login to the registry
	// var loginResponse registry.AuthenticateOKBody
	// loginResponse, err = docker.RegistryLogin(dockerContext, dockerTypes.AuthConfig{
	// 	ServerAddress: dockerRegistryURL,
	// })
	// if err != nil {
	// 	return nil, nil, err
	// }
	// logger.Default.Verbose.Println(fmt.Sprintf("Login - %s: %s", loginResponse.Status, loginResponse.IdentityToken))
	// logger.Default.Verbose.Println(fmt.Sprintf("Created Docker client for registry: %s", dockerRegistryURL))

	return dockerContext, docker, nil
}

// GetImageName creates a new image name based on the Docker repository, package name and resolved package version.
func GetImageName(dockerRepo *string, packageName string, resolvedPackageVersion string) (string, error) {
	var dockerRepoWithSlash string
	if dockerRepo != nil {
		dockerRepoWithSlash = (*dockerRepo) + "/"
	}
	var imageName = fmt.Sprintf("%s%s:%s", dockerRepoWithSlash, packageName, resolvedPackageVersion)

	return imageName, nil
}

// GetDockerfile returns the string contents of a Dockerfile
func GetDockerfile(kpmHomeDir string, packageFullName string) string {
	// Create Dockerfile string
	var dockerfile = strings.TrimSpace(fmt.Sprintf(`
FROM scratch
COPY %s/ /%s
`, packageFullName, packageFullName))

	logger.Default.Verbose.Println(fmt.Sprintf("Generated Dockerfile:\n%s", dockerfile))

	return dockerfile
}

// BuildImage builds a new docker image by making a call to the Docker daemon.
func BuildImage(dockerContext context.Context, docker *client.Client, dockerfile string, dirToCopy string, imageName string) error {
	var err error

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
	response, err = docker.ImageBuild(dockerContext, buildRequestBytes, buildOptions)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Print output
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	logger.Default.Verbose.Println(fmt.Sprintf("Docker daemon's reported OS: %s", response.OSType))
	err = jsonmessage.DisplayJSONMessagesStream(response.Body, os.Stderr, termFd, isTerm, nil)
	if err != nil {
		return err
	}

	return nil
}

// PushImage pushes a Docker image to a remote Docker registry.
func PushImage(dockerContext context.Context, docker *client.Client, imageName string) error {
	// var err error

	// var pushOpts = dockerTypes.ImagePushOptions{}

	// // Push image
	// var pushResponse io.ReadCloser
	// pushResponse, err = docker.ImagePush(dockerContext, imageName, pushOpts)
	// if err != nil {
	// 	return err
	// }
	// defer pushResponse.Close()

	// // Print output
	// termFd, isTerm := term.GetFdInfo(os.Stderr)
	// err = jsonmessage.DisplayJSONMessagesStream(pushResponse, os.Stderr, termFd, isTerm, nil)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// DeleteImage deletes a local Docker image.
func DeleteImage(dockerContext context.Context, docker *client.Client, imageName string) error {
	var err error

	// Set options
	var removeOpts = dockerTypes.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}

	// Delete the image
	var removeResponse []dockerTypes.ImageDeleteResponseItem
	removeResponse, err = docker.ImageRemove(dockerContext, imageName, removeOpts)
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

// createTar creates a new in-memory Docker tar file which can be used to build Docker images by
// making a request to the Docker daemon.
func createTar(dockerfileName string, dockerfile string, dirToCopy string) (*bytes.Buffer, error) {
	var err error

	var dirToCopyBaseName = filepath.Base(dirToCopy)

	// Create byte stream
	var byteStream = new(bytes.Buffer)

	var tarFile = archiver.Tar{
		ImplicitTopLevelFolder: false,
		ContinueOnError:        false,
	}

	// Create a new instance of the tar file in memory
	err = tarFile.Create(byteStream)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()

	// Add Dockerfile
	err = tarFile.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo: types.MockFileInfo{
				MockName:    dockerfileName,
				MockSize:    int64(len(dockerfile)),
				MockMode:    os.ModePerm,
				MockModTime: time.Now(),
			},
			CustomName: dockerfileName,
		},
		ReadCloser: ioutil.NopCloser(bytes.NewBufferString(dockerfile)),
	})

	// Walk the package directory and add all files
	err = filepath.Walk(dirToCopy, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Don't add this file to the tar archive if it is a directory or symbolic link
		if !fileInfo.Mode().IsRegular() {
			return nil
		}

		// Get the internal file path
		var internalFilePath string
		internalFilePath, err = filepath.Rel(dirToCopy, filePath)
		if err != nil {
			return err
		}

		// Prepend the full package name to the path so the files don't all get copied to the root
		internalFilePath = filepath.Join(dirToCopyBaseName, internalFilePath)

		// Replace all backslashes with forward slashes since Docker uses Unix file paths
		internalFilePath = strings.Replace(internalFilePath, "\\", "/", -1)

		logger.Default.Verbose.Println(fmt.Sprintf("Include: %s", internalFilePath))

		// Open the file for reading
		var fileReader *os.File
		fileReader, err = os.Open(filePath)
		if err != nil {
			return err
		}

		// Write the file to the archive
		err = tarFile.Write(archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   fileInfo,
				CustomName: internalFilePath,
			},
			ReadCloser: fileReader,
		})
		fileReader.Close()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Default.Verbose.Println(fmt.Sprintf("Tar file size: %d bytes", len(byteStream.Bytes())))

	var tempOutputDir = filepath.Join(os.TempDir(), "kpm")
	err = os.MkdirAll(tempOutputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	var tempOutputPath = filepath.Join(tempOutputDir, fmt.Sprintf("%s.tar", dirToCopyBaseName))
	err = ioutil.WriteFile(tempOutputPath, byteStream.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}
	logger.Default.Verbose.Println(fmt.Sprintf("Wrote temp file to: %s", tempOutputPath))

	return byteStream, nil
}
