package subcommands

import (
	"fmt"

	"./utils/constants"
	"./utils/docker"
	"./utils/files"
	"./utils/logger"
	"./utils/validation"
)

// PullCmd pulls a template package from a Docker registry to the local filesystem.
func PullCmd(dockerRegistryArg *string, packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string) error {
	var err error

	// Resolve KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDirPath)
	if err != nil {
		return err
	}

	// Get Docker registry URL
	var dockerRegistry = validation.GetStringOrDefault(dockerRegistryArg, docker.DefaultDockerRegistry)

	// Get package name
	var packageName string
	packageName, err = validation.GetStringOrError(packageNameArg, "packageName")
	if err != nil {
		return err
	}

	// Get version
	var resolvedPackageVersion string
	resolvedPackageVersion, err = validation.GetStringOrError(packageVersionArg, "packageVersion")
	if err != nil {
		return err
	}

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Validate package version
	err = validation.ValidatePackageVersion(resolvedPackageVersion, false)
	if err != nil {
		return err
	}

	// Get the package repository directory
	var packageRepositoryDir = constants.GetPackageRepositoryDirPath(kpmHomeDir)

	// Get the image name
	var imageName = docker.GetImageName(dockerRegistry, packageName, resolvedPackageVersion)

	// Pull the Docker image
	err = docker.PullImage(imageName)
	if err != nil {
		return err
	}

	// Delete the local image after we're done
	defer func() {
		var deleteErr = docker.DeleteImage(imageName)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to delete image:\n%s\n%s", deleteErr, err)
			}
		}
	}()

	// Get the package's full name
	var packageFullName = constants.GetPackageFullName(packageName, resolvedPackageVersion)

	// Get the package directory
	var packageDir = constants.GetPackageDirPath(packageRepositoryDir, packageFullName)

	// Extract Docker image contents into the local package repository
	err = docker.ExtractImageContents(imageName, packageDir)
	if err != nil {
		logger.Default.Error.Println("Failed to copy Docker image contents")
		return err
	}

	return err
}
