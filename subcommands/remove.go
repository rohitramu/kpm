package subcommands

import (
	"fmt"

	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// RemoveCmd removes a template package from the local KPM repository.
func RemoveCmd(packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string, userHasConfirmedArg *bool) error {
	var err error

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDir)
	if err != nil {
		return err
	}

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
			return fmt.Errorf("Could not find package '%s' in the local KPM repository. %s", packageName, err)
		}
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return err
	}

	// Resolve package path
	var packageFullName = constants.GetPackageFullName(packageName, packageVersion)
	var packageDir = constants.GetPackageDir(kpmHomeDir, packageFullName)

	var userHasConfirmed bool = validation.GetBoolOrDefault(userHasConfirmedArg, false)
	if err = files.DeleteDirIfExists(packageDir, "template package", userHasConfirmed); err != nil {
		return err
	}

	return nil
}
