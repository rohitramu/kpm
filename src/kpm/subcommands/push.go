package subcommands

import (
	"./common"
	"./utils/docker"
	"./utils/files"
	"./utils/logger"
	"./utils/validation"
)

// PushCmd pushes the template package to a Docker registry.
func PushCmd(kpmHomeDirPathArg *string, dockerRegistryURLArg *string, dockerNamespaceArg *string, packageNameArg *string, packageVersionArg *string) error {
	var err error

	// Resolve KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, common.GetDefaultKpmHomeDirPath)
	if err != nil {
		return err
	}

	// Get Docker registry URL
	var dockerRegistryURL = validation.GetStringOrDefault(dockerRegistryURLArg, docker.DefaultDockerRegistryURL)

	// Get Docker image namespace
	var dockerNamespace = dockerNamespaceArg

	// Validate name
	var packageName string
	packageName, err = validation.GetStringOrError(packageNameArg, "packageName")
	if err != nil {
		return err
	}

	// Validate version
	var wildcardPackageVersion = validation.GetStringOrDefault(packageVersionArg, "*")

	// Resolve the package version
	var resolvedPackageVersion string
	resolvedPackageVersion, err = common.ResolvePackageVersion(kpmHomeDir, packageName, wildcardPackageVersion)
	if err != nil {
		return err
	}

	// Get the full name of the package
	var packageFullName = common.GetPackageFullName(packageName, resolvedPackageVersion)

	// Get the package directory
	var packageDir = common.GetPackageDirPath(common.GetPackageRepositoryDirPath(kpmHomeDir), packageFullName)

	// Create the Dockerfile
	var dockerfile = docker.GetDockerfile(kpmHomeDir, packageFullName)

	// Get a Docker client
	dockerContext, dockerClient, err := docker.GetClient(dockerRegistryURL)
	if err != nil {
		return err
	}

	// Create the image name
	var imageName string
	imageName, err = docker.GetImageName(dockerNamespace, packageName, resolvedPackageVersion)
	if err != nil {
		return err
	}

	// Build the Docker image
	err = docker.BuildImage(dockerContext, dockerClient, dockerfile, packageDir, imageName)
	if err != nil {
		return err
	}

	// Push the docker image
	err = docker.PushImage(dockerContext, dockerClient, imageName)
	if err != nil {
		// Don't error out here, since we still want to try to clean up the created Docker images
		logger.Default.Warning.Println(err)
	}

	// Delete the local image
	err = docker.DeleteImage(dockerContext, dockerClient, imageName)
	if err != nil {
		return err
	}

	return nil
}
