package template_package

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/validation"
)

// PackageNamesAndVersions is the mapping of package names to its list of versions.
type PackageNamesAndVersions map[string][]string

// GetOutputFriendlyName returns the friendly name given an output name and a package's full name.
func GetOutputFriendlyName(outputName string, packageFullName string) string {
	if outputName == packageFullName {
		return packageFullName
	}

	return fmt.Sprintf("%s (%s)", outputName, packageFullName)
}

// GetHighestPackageVersion returns the highest available package version found in the local KPM repository.
func GetHighestPackageVersion(kpmHomeDir string, packageName string) (string, error) {
	var err error

	// Get all available package names and versions
	var availablePackagesAndVersions PackageNamesAndVersions
	availablePackagesAndVersions, err = GetAvailablePackagesAndVersions(kpmHomeDir)
	if err != nil {
		return "", err
	}

	// For each version, resolve the version number
	var availableVersions []string
	var found bool
	availableVersions, found = availablePackagesAndVersions[packageName]
	if !found {
		return "", fmt.Errorf("unable to find template package \"%s\" in local KPM package repository: %s", packageName, kpmHomeDir)
	}

	// Make sure the array is not empty
	if len(availableVersions) == 0 {
		return "", fmt.Errorf("no versions of the template package \"%s\" were found in the local KPM repository: %s", packageName, kpmHomeDir)
	}

	// Get the highest available version
	var highestVersion *string
	for _, currentVersion := range availableVersions {
		// Keep replacing the current version if we found a higher matching version until we get to the end of the matched list
		if highestVersion == nil || currentVersion > *highestVersion {
			highestVersion = &currentVersion
		}
	}

	// This value will never be null since we already checked that the array is not empty
	var result = *highestVersion

	return result, nil
}

// GetAvailablePackagesAndVersions retrieves the list of available packages and their versions.
func GetAvailablePackagesAndVersions(kpmHomeDir string) (PackageNamesAndVersions, error) {
	var err error

	// Get the full list of package names
	var packages []string
	packages, err = GetPackageFullNamesFromLocalRepository(kpmHomeDir)
	if err != nil {
		return nil, err
	}

	// Iterate over the package full names
	var availablePackagesAndVersions = PackageNamesAndVersions{}
	for _, currentPackage := range packages {
		// Extract name and version
		currentPackageName, currentPackageVersion, err := validation.ExtractNameAndVersionFromPackageFullName(currentPackage)
		if err != nil {
			return nil, err
		}

		// If an entry doesn't exist yet for this package version, create it
		var versionsForPackage, ok = availablePackagesAndVersions[currentPackageName]
		if !ok {
			versionsForPackage = []string{}
		}

		// Add the current version to the list of versions for the current package
		availablePackagesAndVersions[currentPackageName] = append(versionsForPackage, currentPackageVersion)
	}

	return availablePackagesAndVersions, nil
}
