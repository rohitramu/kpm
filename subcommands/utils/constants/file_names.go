package constants

import (
	"fmt"
	"path/filepath"
)

// InterfaceFileName is the interface file's name
const InterfaceFileName = "interface.yaml"

// PackageInfoFileName is the package info file's name
const PackageInfoFileName = "package.yaml"

// ParametersFileName is the parameters file's name
const ParametersFileName = "parameters.yaml"

// GetDefaultParametersFile returns the path of the default parameters file in a template package.
func GetDefaultParametersFile(packageDir string) string {
	var parametersFilePath = filepath.Join(packageDir, ParametersFileName)

	return parametersFilePath
}

// GetInterfaceFile returns the path of the interface file in a template package.
func GetInterfaceFile(packageDir string) string {
	var interfaceFilePath = filepath.Join(packageDir, InterfaceFileName)

	return interfaceFilePath
}

// GetPackageInfoFile returns the path of the package information file in a template package.
func GetPackageInfoFile(packageDir string) string {
	var packageInfoFilePath = filepath.Join(packageDir, PackageInfoFileName)

	return packageInfoFilePath
}

// GetPackageFullName returns the full package name with version.
func GetPackageFullName(packageName string, packageVersion string) string {
	return fmt.Sprintf("%s-%s", packageName, packageVersion)
}
