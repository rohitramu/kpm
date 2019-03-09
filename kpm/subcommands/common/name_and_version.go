package common

import (
	"fmt"
	"strings"

	"../utils/constants"
	"../utils/log"
	"../utils/validation"
)

type packageNamesAndVersions map[string][]string

// ResolvePackageVersion returns the highest available package version found in the local KPM repository, which is compatible with a given wildcard package version.
func ResolvePackageVersion(kpmHomeDir string, packageName string, wildcardPackageVersion string) (string, error) {
	var err error

	var packagesDir = constants.GetPackageRepositoryDirPath(kpmHomeDir)

	var resolvedPackageVersion string
	if !strings.Contains(wildcardPackageVersion, "*") {
		// Since this version doesn't have any wildcards, just use it as-is
		resolvedPackageVersion = wildcardPackageVersion
	} else {
		// Get all available package names and versions
		var availablePackagesAndVersions packageNamesAndVersions
		availablePackagesAndVersions, err = getAvailablePackagesAndVersions(packagesDir)
		if err != nil {
			return "", err
		}

		// For each version, resolve the version number
		if availableVersions, found := availablePackagesAndVersions[packageName]; found {
			// Resolve wildcards if required
			resolvedPackageVersion, err = resolveVersionNumber(wildcardPackageVersion, availableVersions)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("Unable to find template package \"%s\" (version: %s) in local KPM package repository: %s", packageName, wildcardPackageVersion, packagesDir)
		}
	}

	return resolvedPackageVersion, nil
}

// getAvailablePackagesAndVersions retrieves the list of available packages and their versions.
func getAvailablePackagesAndVersions(packagesDir string) (packageNamesAndVersions, error) {
	var err error

	// Get the full list of package names
	var packagesList []string
	packagesList, err = GetPackageNamesFromLocalRepository(packagesDir)
	if err != nil {
		return nil, err
	}

	// Iterate over the package full names
	var availablePackagesAndVersions = packageNamesAndVersions{}
	for _, currentPackage := range packagesList {
		// Extract name and version
		currentPackageName, currentPackageVersion, err := validation.ExtractNameAndVersionFromFullPackageName(currentPackage)
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

func resolveVersionNumber(wildcardVersion string, availableVersions []string) (string, error) {
	// Make sure the version is valid
	if err := validation.ValidatePackageVersion(wildcardVersion, true); err != nil {
		log.Panic(err)
	}

	// If the version has a wildcard, get the version up until (and not including) the wildcard character
	var versionWithoutWildcards = wildcardVersion
	if wildcardIndex := strings.IndexRune(wildcardVersion, '*'); wildcardIndex >= 0 {
		versionWithoutWildcards = wildcardVersion[:wildcardIndex]
	}

	// Get the highest available version as specified by the wildcard
	var highestVersion *string
	for _, currentVersion := range availableVersions {
		// Keep replacing the current version if we found a higher matching version until we get to the end of the matched list
		if strings.HasPrefix(currentVersion, versionWithoutWildcards) && (highestVersion == nil || currentVersion > *highestVersion) {
			highestVersion = &currentVersion
		}
	}

	if highestVersion == nil {
		return "", fmt.Errorf("Unable to find a compatible version to resolve: %s", wildcardVersion)
	}

	return *highestVersion, nil
}
