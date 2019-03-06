package common

import (
	"path/filepath"

	"../utils/constants"
	"../utils/files"
)

// GetDefaultKpmHomeDirPath returns the default location of the KPM home directory.
func GetDefaultKpmHomeDirPath() (string, error) {
	var err error

	var userHomeDir string
	userHomeDir, err = files.GetUserHomeDir()
	if err != nil {
		return "", err
	}

	var result = filepath.Join(userHomeDir, constants.KpmHomeDirName)

	return result, nil
}

// GetPackageRepositoryDirPath returns the location of the local package repository.
func GetPackageRepositoryDirPath(kpmHomeDir string) string {
	var packageRepositoryDir = filepath.Join(kpmHomeDir, constants.PackageRepositoryDirName)

	return packageRepositoryDir
}

// GetOutputDirPath returns the path of the output directory for generated files.
func GetOutputDirPath(rootOutputDir string, outputName string) string {
	var outputDirPath = filepath.Join(rootOutputDir, constants.GeneratedDirName, filepath.Base(outputName))

	return outputDirPath
}

// GetPackageDirPath returns the location of a template package in the KPM home directory.
func GetPackageDirPath(packageRepository string, packageFullName string) string {
	var packageDir = filepath.Join(packageRepository, packageFullName)

	return packageDir
}

// GetDependenciesDirPath returns the path of the dependency definition directory in a template package.
func GetDependenciesDirPath(packageDirPath string) string {
	var dependenciesDirPath = filepath.Join(packageDirPath, constants.DependenciesDirName)

	return dependenciesDirPath
}

// GetHelpersDirPath returns the path of the "helpers" directory in a template package.
func GetHelpersDirPath(packageDirPath string) string {
	var helpersDirPath = filepath.Join(packageDirPath, constants.HelpersDirName)

	return helpersDirPath
}

// GetTemplatesDirPath returns the path of the templates directory in a template package.
func GetTemplatesDirPath(packageDirPath string) string {
	var templatesDirPath = filepath.Join(packageDirPath, constants.TemplatesDirName)

	return templatesDirPath
}

// GetDefaultParametersFilePath returns the path of the default parameters file in a template package.
func GetDefaultParametersFilePath(packageDirPath string) string {
	var parametersFilePath = filepath.Join(packageDirPath, constants.ParametersFileName)

	return parametersFilePath
}

// GetInterfaceFilePath returns the path of the interface file in a template package.
func GetInterfaceFilePath(packageDirPath string) string {
	var interfaceFilePath = filepath.Join(packageDirPath, constants.InterfaceFileName)

	return interfaceFilePath
}

// GetPackageInfoFilePath returns the path of the package information file in a template package.
func GetPackageInfoFilePath(packageDirPath string) string {
	var packageInfoFilePath = filepath.Join(packageDirPath, constants.PackageInfoFileName)

	return packageInfoFilePath
}
