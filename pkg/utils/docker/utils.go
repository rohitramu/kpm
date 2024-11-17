package docker

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

// DefaultDockerRegistry is the default registry to use (Docker Hub).
const DefaultDockerRegistry = "docker.io"

// DockerfileRootDir is the root directory to use when building or copying from a Docker image.
const DockerfileRootDir = ".kpm"

func GetImageNameWithoutTag(dockerRegistry string, packageName string) string {
	imageName := packageName
	if dockerRegistry != DefaultDockerRegistry {
		imageName = fmt.Sprintf("%s/%s", dockerRegistry, imageName)
	}

	return imageName
}

// GetImageName creates a new image name based on the Docker repository, package name and resolved package version.
func GetImageName(dockerRegistry string, packageName string, resolvedPackageVersion string) string {
	imageName := fmt.Sprintf("%s:%s", GetImageNameWithoutTag(dockerRegistry, packageName), resolvedPackageVersion)

	return imageName
}

// GetDockerfilePath returns the path of the Dockerfile to use.
func GetDockerfilePath(kpmHomeDir string) string {
	var dockerfilePath = filepath.Join(kpmHomeDir, "Dockerfile")

	// If the file doesn't exist, create it
	if err := files.FileExists(dockerfilePath, "Dockerfile"); err != nil {
		// Create Dockerfile string
		var dockerfile = fmt.Sprintf(`
FROM scratch
COPY ./ /%s
CMD [""]
`, DockerfileRootDir)
		dockerfile = strings.TrimSpace(dockerfile)

		// Write to file
		err = files.CreateFile(dockerfilePath, "dockerfile", dockerfile)
		if err != nil {
			log.Panicf("Failed to create dockerfile: %s", err)
		}

		log.Debugf("Generated Dockerfile:\n%s", dockerfile)
	}

	return dockerfilePath
}
