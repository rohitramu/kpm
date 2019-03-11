package constants

import (
	"path/filepath"

	"github.com/rohitramu/kpm/subcommands/utils/files"
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

// GetDefaultKpmHomeDir returns the default location of the KPM home directory.
func GetDefaultKpmHomeDir() (string, error) {
	var err error

	var userHomeDir string
	userHomeDir, err = files.GetUserHomeDir()
	if err != nil {
		return "", err
	}

	var result = filepath.Join(userHomeDir, KpmHomeDirName)

	return result, nil
}

// GetPackageRepositoryDir returns the location of the local package repository.
func GetPackageRepositoryDir(kpmHomeDir string) string {
	var packageRepositoryDir = filepath.Join(kpmHomeDir, PackageRepositoryDirName)

	return packageRepositoryDir
}

// GetDefaultOutputDir returns the default path of the root directory for generated files.
func GetDefaultOutputDir(outputParentDir string) string {
	var outputDirPath = filepath.Join(outputParentDir, GeneratedDirName)

	return outputDirPath
}

// GetPackageDir returns the location of a template package in the KPM home directory.
func GetPackageDir(kpmHomeDir string, packageFullName string) string {
	var packageDir = filepath.Join(GetPackageRepositoryDir(kpmHomeDir), packageFullName)

	return packageDir
}

// GetDependenciesDir returns the path of the dependency definition directory in a template package.
func GetDependenciesDir(packageDir string) string {
	var dependenciesDirPath = filepath.Join(packageDir, DependenciesDirName)

	return dependenciesDirPath
}

// GetHelpersDir returns the path of the "helpers" directory in a template package.
func GetHelpersDir(packageDir string) string {
	var helpersDirPath = filepath.Join(packageDir, HelpersDirName)

	return helpersDirPath
}

// GetTemplatesDir returns the path of the templates directory in a template package.
func GetTemplatesDir(packageDir string) string {
	var templatesDirPath = filepath.Join(packageDir, TemplatesDirName)

	return templatesDirPath
}
