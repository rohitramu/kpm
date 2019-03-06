package subcommands

import (
	"fmt"

	"./common"
	"./utils/docker"
	"./utils/files"
	"./utils/logger"
	"./utils/validation"
)

// PullCmd pulls a template package from a Docker registry to the local filesystem.
func PullCmd(kpmHomeDirPathArg *string, dockerRegistryURLArg *string, dockerNamespaceArg *string, packageNameArg *string, packageVersionArg *string) error {
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

	// Create the image name
	var imageName string
	imageName, err = docker.GetImageName(dockerNamespace, packageName, resolvedPackageVersion)
	if err != nil {
		return err
	}

	// Pull the Docker image
	logger.Default.Verbose.Println(fmt.Sprintf("Pulling Docker image \"%s\" from: %s", imageName, dockerRegistryURL))
	err = docker.PullImage(dockerRegistryURL, imageName)
	if err != nil {
		return err
	}

	// Delete the local image after we're done
	defer func() {
		logger.Default.Verbose.Println(fmt.Sprintf("Deleting image: %s", imageName))
		var deleteErr = docker.DeleteImage(imageName)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to delete image:\n%s\n%s", deleteErr, err)
			}
		}
	}()

	// Copy Docker image contents into the local package repository
	logger.Default.Verbose.Println(fmt.Sprintf("Copying image contents to: %s", packageDir))
	err = docker.CopyImageContents(imageName, packageFullName, packageDir)
	if err != nil {
		logger.Default.Error.Println("Failed to copy Docker image contents")
		return err
	}

	return nil
}
