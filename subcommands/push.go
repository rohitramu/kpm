package subcommands

import (
	"fmt"

	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/docker"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/types"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// PushCmd pushes the template package to a Docker registry.
func PushCmd(dockerRegistryArg *string, packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string) error {
	var err error

	// Resolve KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDir)
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

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Get package version
	var packageVersion string
	packageVersion, err = validation.GetStringOrError(packageVersionArg, "packageVersion")
	if err != nil {
		// Since the package version was not provided, check the local repository for the highest version
		if packageVersion, err = common.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
			return err
		}
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return err
	}

	// Get the package full name
	var packageFullName = constants.GetPackageFullName(packageName, packageVersion)

	// Get the package directory
	var packageDir = constants.GetPackageDir(kpmHomeDir, packageFullName)

	log.Info("Validating package: %s", packageDir)

	// Get package information
	_, err = common.GetPackageInfo(kpmHomeDir, packageDir)
	if err != nil {
		return err
	}

	// Get the default parameters
	var packageParameters *types.GenericMap
	packageParameters, err = common.GetPackageParameters(constants.GetDefaultParametersFile(packageDir))
	if err != nil {
		return err
	}

	// Ensure that the dependency tree can be calculated without errors using the default parameters
	var outputName = constants.GetDefaultOutputName(packageName, packageVersion)
	_, err = common.GetDependencyTree(kpmHomeDir, packageName, packageVersion, dockerRegistry, outputName, packageParameters)
	if err != nil {
		return err
	}

	// Create the image name
	var imageName = docker.GetImageName(dockerRegistry, packageName, packageVersion)

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
