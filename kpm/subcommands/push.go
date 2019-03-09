package subcommands

import (
	"fmt"

	"./common"
	"./utils/constants"
	"./utils/docker"
	"./utils/files"
	"./utils/validation"
)

// PushCmd pushes the template package to a Docker registry.
func PushCmd(dockerRegistryArg *string, packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string) error {
	var err error

	// Resolve KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDirPath)
	if err != nil {
		return err
	}

	// Get Docker registry name
	var dockerRegistry = validation.GetStringOrDefault(dockerRegistryArg, docker.DefaultDockerRegistry)

	// Get package name
	var packageName string
	packageName, err = validation.GetStringOrError(packageNameArg, "packageName")
	if err != nil {
		return err
	}

	// Get version
	var wildcardPackageVersion = validation.GetStringOrDefault(packageVersionArg, "*")

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Validate package version
	err = validation.ValidatePackageVersion(wildcardPackageVersion, true)
	if err != nil {
		return err
	}

	// Resolve the package version
	var resolvedPackageVersion string
	resolvedPackageVersion, err = common.ResolvePackageVersion(kpmHomeDir, packageName, wildcardPackageVersion)
	if err != nil {
		return err
	}

	// Get the package repository directory
	var packageRepositoryDir = constants.GetPackageRepositoryDirPath(kpmHomeDir)

	// Get the package full name
	var packageFullName = constants.GetPackageFullName(packageName, resolvedPackageVersion)

	// Get the package directory
	var packageDir = constants.GetPackageDirPath(packageRepositoryDir, packageFullName)

	// Validate the package
	_, err = common.GetPackageInfo(packageDir)
	if err != nil {
		return err
	}

	// Create the image name
	var imageName = docker.GetImageName(dockerRegistry, packageName, resolvedPackageVersion)

	// Create the Dockerfile
	var dockerfilePath = docker.GetDockerfilePath(kpmHomeDir)

	// Build the Docker image
	err = docker.BuildImage(imageName, dockerfilePath, packageDir)
	if err != nil {
		return err
	}

	// Delete the local image after we're done
	defer func() {
		var deleteErr = docker.DeleteImage(imageName)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to clean up image: %s\n%s\n%s", imageName, deleteErr, err)
			} else {
				err = deleteErr
			}
		}
	}()

	// Push the Docker image
	err = docker.PushImage(imageName)
	if err != nil {
		return err
	}

	return nil
}
