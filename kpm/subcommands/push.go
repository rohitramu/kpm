package subcommands

import (
	"fmt"

	"./common"
	"./utils/docker"
	"./utils/files"
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

	// Create the image name
	var imageName string
	imageName, err = docker.GetImageName(dockerNamespace, packageName, resolvedPackageVersion)
	if err != nil {
		return err
	}

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

	// Push the Docker image, but don't error out here if it fails (we still need to clean up the image)
	err = docker.PushImage(dockerRegistryURL, imageName)
	if err != nil {
		return err
	}

	return nil
}
