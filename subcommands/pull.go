package subcommands

import (
	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/docker"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// PullCmd pulls a template package from a Docker registry to the local filesystem.
func PullCmd(packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string, dockerRegistryArg *string) error {
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

	// Pull package
	err = common.PullPackage(kpmHomeDir, dockerRegistry, packageName, packageVersion)
	if err != nil {
		return err
	}

	return nil
}
