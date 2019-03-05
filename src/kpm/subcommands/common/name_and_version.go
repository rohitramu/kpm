package common

import (
	"fmt"
	"io/ioutil"
	"strings"

	"../utils/logger"
	"../utils/validation"
)

// GetPackageFullName returns the full package name with version.
func GetPackageFullName(packageName string, resolvedPackageVersion string) string {
	return fmt.Sprintf("%s-%s", packageName, resolvedPackageVersion)
}

// ResolvePackageVersion returns the highest available package version that is compatible with a given wildcard package version.
func ResolvePackageVersion(kpmHomeDir string, packageName string, wildcardPackageVersion string) (string, error) {
	var err error

	var packagesDir = GetPackageRepositoryDirPath(kpmHomeDir)

	var resolvedPackageVersion string
	if !strings.Contains(wildcardPackageVersion, "*") {
		// Since this version doesn't have any wildcards, just use it as-is
		resolvedPackageVersion = wildcardPackageVersion
	} else {
		// Get the names of all available versions of the package
		var availablePackagesAndVersions = getAvailablePackagesAndVersions(packagesDir)
		if availableVersions, ok := availablePackagesAndVersions[packageName]; ok {
			// Resolve wildcards if required
			resolvedPackageVersion, err = resolveVersionNumber(wildcardPackageVersion, availableVersions)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("Unable to find template package in local KPM package repository: %s", packagesDir)
		}
	}

	return resolvedPackageVersion, nil
}

// getAvailablePackagesAndVersions retrieves the list of available packages and their versions.
func getAvailablePackagesAndVersions(packagesDir string) map[string][]string {
	var availablePackagesAndVersions = map[string][]string{}
	if files, err := ioutil.ReadDir(packagesDir); err != nil {
		logger.Default.Error.Panicln(err)
	} else {
		for _, file := range files {
			var fileName = file.Name()

			// Ensure that we are looking at a directory
			if file.IsDir() {
				currentPackageName, currentPackageVersion, err := validation.ExtractNameAndVersionFromFullPackageName(fileName)
				if err != nil {
					logger.Default.Verbose.Println(fmt.Sprintf("Found non-package directory \"%s\": %s", fileName, err))
				} else {
					// If an entry doesn't exist yet for this package version, create it
					var versionsForPackage, ok = availablePackagesAndVersions[currentPackageName]
					if !ok {
						versionsForPackage = []string{}
					}

					// Add the current version to the list of versions for the current package
					availablePackagesAndVersions[currentPackageName] = append(versionsForPackage, currentPackageVersion)
				}
			}
		}
	}

	return availablePackagesAndVersions
}

func resolveVersionNumber(wildcardVersion string, availableVersions []string) (string, error) {
	// Make sure the version is valid
	if err := validation.ValidatePackageVersion(wildcardVersion, true); err != nil {
		logger.Default.Error.Panicln(err)
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
