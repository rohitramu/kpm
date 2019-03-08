package docker

import (
	"fmt"
	"strings"

	"../logger"
)

// DefaultDockerRegistryURL is the default registry URL (Docker Hub).
const DefaultDockerRegistryURL = "https://index.docker.io/v1/"

// DockerTarFileRootDir is the root directory to use when building or copying from a Docker image.
const DockerTarFileRootDir = ".kpm"

// GetImageName creates a new image name based on the Docker repository, package name and resolved package version.
func GetImageName(packageName string, resolvedPackageVersion string) string {
	var imageName = packageName + ":" + resolvedPackageVersion

	return imageName
}

// GetDockerfile returns the string contents of a Dockerfile
func GetDockerfile() string {
	// Create Dockerfile string
	var dockerfile = strings.TrimSpace(fmt.Sprintf(`
FROM scratch
COPY %s/ /%s
`, DockerTarFileRootDir, DockerTarFileRootDir))

	logger.Default.Verbose.Println(fmt.Sprintf("Generated Dockerfile:\n%s", dockerfile))

	return dockerfile
}
