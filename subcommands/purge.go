package subcommands

import (
	"fmt"
	"os"

	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/user_prompts"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// PurgeCmd removes all versions of a template package from the local KPM repository.
func PurgeCmd(packageNameArg *string, userHasConfirmedArg *bool, kpmHomeDirPathArg *string) error {
	var err error
	var ok bool

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDir)
	if err != nil {
		return err
	}

	// Get package name
	var packageName = validation.GetStringOrDefault(packageNameArg, "")

	// Find all packages and versions
	var packages common.PackageNamesAndVersions
	packages, err = common.GetAvailablePackagesAndVersions(kpmHomeDir)
	if err != nil {
		return fmt.Errorf("Failed to retrieve the list of available packages in the local KPM repository: %s", err)
	}

	// Create function to remove all versions of a package
	var removeAllVersions = func(currentPackageName string) error {
		// Get the versions of this package
		var versions []string
		versions, ok = packages[currentPackageName]
		if !ok {
			return fmt.Errorf("Failed to find package in the local KPM repository: %s", currentPackageName)
		}

		var userHasConfirmed bool = validation.GetBoolOrDefault(userHasConfirmedArg, false)
		if !userHasConfirmed {
			if userHasConfirmed, err = user_prompts.ConfirmWithUser("All versions of package '%s' will be deleted from the local KPM repository.", currentPackageName); err != nil {
				log.Panic("Failed to get user confirmation. \n%s", err)
			}

			if !userHasConfirmed {
				return fmt.Errorf("Purge operation cancelled - user did not confirm the delete action.")
			}
		}

		// Remove the package
		for _, ver := range versions {
			var packageFullName = constants.GetPackageFullName(currentPackageName, ver)
			var packageDir = constants.GetPackageDir(kpmHomeDir, packageFullName)
			err = os.RemoveAll(packageDir)
			if err != nil {
				return fmt.Errorf("Failed to remove package \"%s\" from the local KPM repository: %s", packageFullName, err)
			}
		}

		return nil
	}

	// If we only want to remove all versions of a single package, iterate over just that package's versions
	if packageName != "" {
		// Validate package name
		err = validation.ValidatePackageName(packageName)
		if err != nil {
			return err
		}

		err = removeAllVersions(packageName)
		if err != nil {
			return err
		}
	} else {
		// Purge the local KPM repository of all packages
		for packageName := range packages {
			err = removeAllVersions(packageName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
