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

// GetDefaultParametersFilePath returns the path of the default parameters file in a template package.
func GetDefaultParametersFilePath(packageDirPath string) string {
	var parametersFilePath = filepath.Join(packageDirPath, ParametersFileName)

	return parametersFilePath
}

// GetInterfaceFilePath returns the path of the interface file in a template package.
func GetInterfaceFilePath(packageDirPath string) string {
	var interfaceFilePath = filepath.Join(packageDirPath, InterfaceFileName)

	return interfaceFilePath
}

// GetPackageInfoFilePath returns the path of the package information file in a template package.
func GetPackageInfoFilePath(packageDirPath string) string {
	var packageInfoFilePath = filepath.Join(packageDirPath, PackageInfoFileName)

	return packageInfoFilePath
}

// GetPackageFullName returns the full package name with version.
func GetPackageFullName(packageName string, resolvedPackageVersion string) string {
	return fmt.Sprintf("%s-%s", packageName, resolvedPackageVersion)
}
