package constants

import (
	"path/filepath"

	"../files"
)

// DependenciesDirName is the name of the directory where package dependencies are defined
const DependenciesDirName = "dependencies"

// TemplatesDirName is the name of the directory where package templates are defined
const TemplatesDirName = "templates"

// HelpersDirName is the name of the directory where package helpers are defined
const HelpersDirName = "helpers"

// KpmHomeDirName is the name of the directory that acts as the home directory for KPM
const KpmHomeDirName = ".kpm"

// PackageRepositoryDirName is the name of the directory that contains packages available for use
const PackageRepositoryDirName = "packages"

// GeneratedDirName is the name of the directory in which generated output is written
const GeneratedDirName = ".kpm_generated"

// GetDefaultKpmHomeDirPath returns the default location of the KPM home directory.
func GetDefaultKpmHomeDirPath() (string, error) {
	var err error

	var userHomeDir string
	userHomeDir, err = files.GetUserHomeDir()
	if err != nil {
		return "", err
	}

	var result = filepath.Join(userHomeDir, KpmHomeDirName)

	return result, nil
}

// GetPackageRepositoryDirPath returns the location of the local package repository.
func GetPackageRepositoryDirPath(kpmHomeDir string) string {
	var packageRepositoryDir = filepath.Join(kpmHomeDir, PackageRepositoryDirName)

	return packageRepositoryDir
}

// GetOutputDirPath returns the path of the output directory for generated files.
func GetOutputDirPath(rootOutputDir string, outputName string) string {
	var outputDirPath = filepath.Join(rootOutputDir, GeneratedDirName, filepath.FromSlash(outputName))

	return outputDirPath
}

// GetPackageDirPath returns the location of a template package in the KPM home directory.
func GetPackageDirPath(packageRepositoryDir string, packageFullName string) string {
	var packageDir = filepath.Join(packageRepositoryDir, packageFullName)

	return packageDir
}

// GetDependenciesDirPath returns the path of the dependency definition directory in a template package.
func GetDependenciesDirPath(packageDirPath string) string {
	var dependenciesDirPath = filepath.Join(packageDirPath, DependenciesDirName)

	return dependenciesDirPath
}

// GetHelpersDirPath returns the path of the "helpers" directory in a template package.
func GetHelpersDirPath(packageDirPath string) string {
	var helpersDirPath = filepath.Join(packageDirPath, HelpersDirName)

	return helpersDirPath
}

// GetTemplatesDirPath returns the path of the templates directory in a template package.
func GetTemplatesDirPath(packageDirPath string) string {
	var templatesDirPath = filepath.Join(packageDirPath, TemplatesDirName)

	return templatesDirPath
}
