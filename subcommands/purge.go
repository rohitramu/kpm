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

// PurgeCmd removes all versions of a template package from the local KPM repository.
func PurgeCmd(packageNameArg *string, allConfirm *bool, kpmHomeDirPathArg *string) error {
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
		return fmt.Errorf("Failed to retrieve the list of available packages: %s", err)
	}

	// Create function to remove all versions of a package
	var removeAllVersions = func(currentPackageName string) error {
		// Get the versions of this package
		var versions []string
		versions, ok = packages[currentPackageName]
		if !ok {
			return fmt.Errorf("Failed to find package in the local KPM repository: %s", currentPackageName)
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

	// If we only want to remove all versions for a single package, iterate over just that package's versions
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
	} else if allConfirm != nil && *allConfirm {
		// If the user has confirmed that they want to remove all packages, purge the local KPM repository
		for packageName := range packages {
			err = removeAllVersions(packageName)
			if err != nil {
				return err
			}
		}
	} else {
		log.Warning("Either provide a package name for which to remove all versions, or confirm that you want to remove all packages from the local KPM repository by providing the \"--%s\" flag", constants.PurgeAllFlagName)
		return nil
	}

	return nil
}
