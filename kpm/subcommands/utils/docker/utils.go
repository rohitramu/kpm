package docker

import (
	"fmt"
	"strings"

	"../logger"
)

// DefaultDockerRegistryURL is the default registry URL (Docker Hub)
const DefaultDockerRegistryURL = "https://index.docker.io/v1/"

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
