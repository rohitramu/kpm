package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
)

// DefaultDockerRegistry is the default registry to use (Docker Hub).
const DefaultDockerRegistry = "docker.io"

// DockerfileRootDir is the root directory to use when building or copying from a Docker image.
const DockerfileRootDir = ".kpm"

// GetImageName creates a new image name based on the Docker repository, package name and resolved package version.
func GetImageName(dockerRegistry string, packageName string, resolvedPackageVersion string) string {
	var imageName = fmt.Sprintf("%s:%s", packageName, resolvedPackageVersion)
	if dockerRegistry != DefaultDockerRegistry {
		imageName = fmt.Sprintf("%s/%s", dockerRegistry, imageName)
	}

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
		err = ioutil.WriteFile(dockerfilePath, []byte(dockerfile), os.ModePerm)
		if err != nil {
			log.Panic("Failed to write dockerfile: %s", err)
		}

		log.Debug("Generated Dockerfile:\n%s", dockerfile)
	}

	return dockerfilePath
}
