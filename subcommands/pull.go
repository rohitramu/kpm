package subcommands

import (
	"fmt"

	"./utils/constants"
	"./utils/docker"
	"./utils/files"
	"./utils/log"
	"./utils/validation"
)

// PullCmd pulls a template package from a Docker registry to the local filesystem.
func PullCmd(dockerRegistryArg *string, packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string) error {
	var err error

	// Resolve KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDir)
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
	var packageVersion string
	packageVersion, err = validation.GetStringOrError(packageVersionArg, "packageVersion")
	if err != nil {
		return err
	}

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return err
	}

	// Get the image name
	var imageName = docker.GetImageName(dockerRegistry, packageName, packageVersion)

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
				err = fmt.Errorf("Failed to delete image: %s\n%s\n%s", imageName, deleteErr, err)
			}
		}
	}()

	// Get the package's full name
	var packageFullName = constants.GetPackageFullName(packageName, packageVersion)

	// Get the package directory
	var packageDir = constants.GetPackageDir(kpmHomeDir, packageFullName)

	// Extract Docker image contents into the local package repository
	err = docker.ExtractImageContents(imageName, packageDir)
	if err != nil {
		log.Error("Failed to copy Docker image contents")
		return err
	}

	return nil
}
