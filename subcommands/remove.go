package subcommands

import (
	"fmt"
	"os"

	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// RemoveCmd removes a template package from the local KPM repository.
func RemoveCmd(packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string) error {
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
			return fmt.Errorf("Package version must be provided if the package does not exist in the local repository: %s", err)
		}
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return err
	}

	// Resolve paths
	var packageFullName = constants.GetPackageFullName(packageName, packageVersion)
	var packageDir = constants.GetPackageDir(kpmHomeDir, packageFullName)

	// Log resolved values
	log.Info("====")
	log.Info("Package name:      %s", packageName)
	log.Info("Package version:   %s", packageVersion)
	log.Info("Package directory: %s", packageDir)
	log.Info("====")

	// Check that the package exists in the local repository
	err = files.DirExists(packageDir, "template package")
	if err != nil {
		// Package doesn't exist, so just log a warning and exit
		log.Warning("Package \"%s\" not found in local repository", packageFullName)
		return nil
	}

	// Delete the package from the local repository
	err = os.RemoveAll(packageDir)
	if err != nil {
		return fmt.Errorf("Failed to remove package from local repository: %s\n%s", packageFullName, err)
	}

	return nil
}
