package template_package

import (
	"fmt"
	"path/filepath"

	"github.com/rohitramu/kpm/pkg/utils/constants"
)

// GetDefaultParametersFile returns the path of the default parameters file in a template package.
func GetDefaultParametersFile(packageDir string) string {
	var parametersFilePath = filepath.Join(packageDir, constants.ParametersFileName)

	return parametersFilePath
}

// GetInterfaceFile returns the path of the interface file in a template package.
func GetInterfaceFile(packageDir string) string {
	var interfaceFilePath = filepath.Join(packageDir, constants.InterfaceFileName)

	return interfaceFilePath
}

// GetPackageInfoFile returns the path of the package information file in a template package.
func GetPackageInfoFile(packageDir string) string {
	var packageInfoFilePath = filepath.Join(packageDir, constants.PackageInfoFileName)

	return packageInfoFilePath
}

// GetPackageFullName returns the full package name with version.
func GetPackageFullName(packageName string, packageVersion string) string {
	return fmt.Sprintf("%s-%s", packageName, packageVersion)
}
