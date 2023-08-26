package template_package

import (
	"path/filepath"

	"github.com/rohitramu/kpm/pkg/utils/constants"
)

// GetRepoPackagesDir returns the location of the local package repository.
func GetRepoPackagesDir(repoDir string) string {
	var packageRepositoryDir = filepath.Join(repoDir, constants.PackageRepositoryDirName)

	return packageRepositoryDir
}

// GetDefaultOutputName returns the default output name when executing a package.
func GetDefaultOutputName(packageName string, packageVersion string) string {
	return GetPackageFullName(packageName, packageVersion)
}

// GetDefaultExportName returns the default name when exporting a package.
func GetDefaultExportName(packageName string, packageVersion string) string {
	return GetPackageFullName(packageName, packageVersion)
}

// GetPackageDir returns the location of a template package in the KPM home directory.
func GetPackageDir(kpmHomeDir string, packageFullName string) string {
	var packageDir = filepath.Join(GetRepoPackagesDir(kpmHomeDir), packageFullName)

	return packageDir
}

// GetDependenciesDir returns the path of the dependency definition directory in a template package.
func GetDependenciesDir(packageDir string) string {
	var dependenciesDirPath = filepath.Join(packageDir, constants.DependenciesDirName)

	return dependenciesDirPath
}

// GetHelpersDir returns the path of the "helpers" directory in a template package.
func GetHelpersDir(packageDir string) string {
	var helpersDirPath = filepath.Join(packageDir, constants.HelpersDirName)

	return helpersDirPath
}

// GetTemplatesDir returns the path of the templates directory in a template package.
func GetTemplatesDir(packageDir string) string {
	var templatesDirPath = filepath.Join(packageDir, constants.TemplatesDirName)

	return templatesDirPath
}
