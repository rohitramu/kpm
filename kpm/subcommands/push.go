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
func PushCmd(dockerRegistryURLArg *string, packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string) error {
	var err error

	// Resolve KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDirPath)
	if err != nil {
		return err
	}

	// Get Docker registry URL
	var dockerRegistryURL = validation.GetStringOrDefault(dockerRegistryURLArg, docker.DefaultDockerRegistryURL)

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
	var imageName = docker.GetImageName(packageName, resolvedPackageVersion)

	// Create the Dockerfile
	var dockerfile = docker.GetDockerfile()

	// Build the Docker image
	err = docker.BuildImage(imageName, dockerfile, packageDir)
	if err != nil {
		return err
	}

	// Delete the local image after we're done
	defer func() {
		var deleteErr = docker.DeleteImage(imageName)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to clean up image:\n%s\n%s", deleteErr, err)
			} else {
				err = deleteErr
			}
		}
	}()

	// Push the Docker image
	err = docker.PushImage(dockerRegistryURL, imageName)
	if err != nil {
		return err
	}

	return nil
}
