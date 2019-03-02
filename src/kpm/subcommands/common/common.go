package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"../utils/constants"
	"../utils/files"
	"../utils/logger"
	"../utils/templates"
	"../utils/types"
	"../utils/validation"
	"../utils/yaml"
)

// PullPackage retrieves a remote template package and makes it available for use.  If a package
// was successfully retrieved, this function returns the retrieved version number.
func PullPackage(packageName string, packageVersion string) (string, error) {
	//TODO: Get list of versions in remote repository

	//TODO: Resolve version to the highest that is compatible with the requested version

	//TODO: Get the template package of the resolved version

	return "", fmt.Errorf("Could not find a compatible version for package: %s", validation.GetFullPackageName(packageName, packageVersion))
}

// GetPackageInput creates the input values for a package by combining the interface, parameters and package info.
func GetPackageInput(parentTemplate *template.Template, packageDirPath string, parametersFilePath string) *types.GenericMap {
	// Get package info
	var packageInfo = GetPackageInfo(packageDirPath)

	// Generate the values by populating the interface template with parameters
	var values = getValuesFromInterface(parentTemplate, packageDirPath, parametersFilePath)

	// Add top-level objects
	var result = types.GenericMap{}
	result[constants.TemplateFieldPackage] = packageInfo
	result[constants.TemplateFieldValues] = values

	return &result
}

// GetPackageInfo returns the package info object for a given package.
func GetPackageInfo(packageDirPath string) *types.PackageInfo {
	var err error

	// Make sure that the package exists
	if fileInfo, err := os.Stat(packageDirPath); err != nil {
		if os.IsNotExist(err) {
			logger.Default.Error.Fatalln(fmt.Sprintf("Package not found in directory: %s", packageDirPath))
		} else {
			logger.Default.Error.Panicln(err)
		}
	} else if !fileInfo.IsDir() {
		logger.Default.Error.Fatalln(fmt.Sprintf("Package path does not point to a directory: %s", packageDirPath))
	} else {
		logger.Default.Verbose.Println(fmt.Sprintf("Found template package in directory: %s", packageDirPath))
	}

	// Check that the package info file exists
	var packageInfoFilePath = filepath.Join(packageDirPath, constants.PackageInfoFileName)
	if fileInfo, err := os.Stat(packageInfoFilePath); err != nil {
		if os.IsNotExist(err) {
			logger.Default.Error.Fatalln(fmt.Sprintf("Package information file does not exist: %s", packageInfoFilePath))
		} else {
			logger.Default.Error.Panicln(err)
		}
	} else if fileInfo.IsDir() {
		logger.Default.Error.Fatalln(fmt.Sprintf("Package path does not point to a file: %s", packageInfoFilePath))
	} else {
		logger.Default.Verbose.Println(fmt.Sprintf("Found package information file: %s", packageInfoFilePath))
	}

	// Get package info file content
	var yamlBytes = files.ReadFileToBytes(packageInfoFilePath)

	// Get package info object from file content
	var packageInfo = new(types.PackageInfo)
	yaml.BytesToObject(yamlBytes, packageInfo)

	// Validate package name
	err = validation.ValidatePackageName(packageInfo.Name)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageInfo.Version, false)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return packageInfo
}

// GetPackageDir returns the location of a template package in the KPM home directory.
func GetPackageDir(kpmHomeDir string, packageName string, packageVersion string) string {
	var err error

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion, true)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Construct packages directory
	var packagesDir = filepath.Join(kpmHomeDir, constants.PackageRepositoryDirName)

	// Resolve the package version
	var resolvedPackageVersion string
	if !strings.Contains(packageVersion, "*") {
		// Since this version doesn't have any wildcards, just use it as-is
		resolvedPackageVersion = packageVersion
	} else {
		// Get the names of all available versions of the package
		var availablePackagesAndVersions = getAvailablePackagesAndVersions(packagesDir)
		if availableVersions, ok := availablePackagesAndVersions[packageName]; ok {
			// Resolve wildcards if required
			resolvedPackageVersion = resolveVersion(packageVersion, availableVersions)
		} else {
			logger.Default.Error.Fatalln(fmt.Sprintf("Unable to find template package in local KPM package repository: %s", packagesDir))
		}
	}

	// Combine the package name and version to get the full package name
	var packageNameWithVersion = validation.GetFullPackageName(packageName, resolvedPackageVersion)

	// Construct the full path to the package directory
	var resolvedPackageDir = filepath.Join(packagesDir, packageNameWithVersion)

	return resolvedPackageDir
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
				currentPackageName, currentPackageVersion, err := validation.GetNameAndVersionFromFullPackageName(fileName)
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

func resolveVersion(wildcardVersion string, availableVersions []string) string {
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
		logger.Default.Error.Fatalln(fmt.Sprintf("Unable to find a compatible version to resolve: %s", wildcardVersion))
	}

	return *highestVersion
}

// getValuesFromInterface creates the values which can be used as input to templates by executing the interface with parameters.
func getValuesFromInterface(parentTemplate *template.Template, packageDirPath string, parametersFilePath string) *types.GenericMap {
	// Create template object from interface file
	var templateName = constants.InterfaceFileName
	var interfaceFilePath = filepath.Join(packageDirPath, templateName)
	var tmpl = templates.GetTemplateFromFile(parentTemplate, templateName, interfaceFilePath)

	// Get parameters file content as bytes
	var parametersFileBytes []byte
	if fileInfo, err := os.Stat(parametersFilePath); err != nil {
		if os.IsNotExist(err) {
			logger.Default.Warning.Println(fmt.Sprintf("Parameters file does not exist: %s", parametersFilePath))
		} else {
			logger.Default.Error.Panicln(err)
		}
	} else if fileInfo.IsDir() {
		logger.Default.Warning.Println(fmt.Sprintf("Parameters file path does not point to a file: %s", parametersFilePath))
	} else {
		logger.Default.Verbose.Println(fmt.Sprintf("Found parameters file: %s", parametersFilePath))
		parametersFileBytes = files.ReadFileToBytes(parametersFilePath)
	}

	// Get parameters
	var parameters = new(types.GenericMap)
	yaml.BytesToObject(parametersFileBytes, parameters)

	// Generate values by applying parameters to interface
	var interfaceBytes = templates.ExecuteTemplate(tmpl, parameters)

	// Get values object from generated values yaml file
	var result = new(types.GenericMap)
	yaml.BytesToObject(interfaceBytes, result)

	return result
}
