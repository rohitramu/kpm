package local_package_repo

import (
	"fmt"
	"os"

	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"golang.org/x/exp/maps"
)

func RemoveAllVersionsOfPackages(kpmHomeDir string, packageNames ...string) error {
	var err error

	if len(packageNames) == 0 {
		return nil
	}

	// Find all packages and versions.
	var packages template_package.PackageNamesAndVersions
	packages, err = template_package.GetAvailablePackagesAndVersions(kpmHomeDir)
	if err != nil {
		return fmt.Errorf("failed to retrieve the list of available packages in the local KPM repository: %s", err)
	}

	return removeAllVersionsOfPackages(kpmHomeDir, packages, packageNames...)
}

func RemoveAllVersionsOfAllPackages(kpmHomeDir string) error {
	var err error

	// Find all packages and versions.
	var packages template_package.PackageNamesAndVersions
	packages, err = template_package.GetAvailablePackagesAndVersions(kpmHomeDir)
	if err != nil {
		return fmt.Errorf("failed to retrieve the list of available packages in the local KPM repository: %s", err)
	}

	return removeAllVersionsOfPackages(kpmHomeDir, packages, maps.Keys(packages)...)
}

func removeAllVersionsOfPackages(
	kpmHomeDir string,
	packages template_package.PackageNamesAndVersions,
	packageNames ...string,
) error {
	var err error

	// Exit early if there are no packages.
	if len(packages) == 0 {
		return nil
	}

	// Delete all versions of packages.
	for _, packageName := range packageNames {
		// Get the versions of this package.
		var versions []string
		var ok bool
		versions, ok = packages[packageName]
		if !ok {
			return fmt.Errorf("failed to find package in the local KPM repository: %s", packageName)
		}

		// Remove the package.
		for _, ver := range versions {
			var packageFullName = template_package.GetPackageFullName(packageName, ver)
			var packageDir = template_package.GetPackageDir(kpmHomeDir, packageFullName)
			err = os.RemoveAll(packageDir)
			if err != nil {
				return fmt.Errorf("failed to remove package \"%s\" from the local KPM repository: %s", packageFullName, err)
			}
		}
	}

	return nil
}
