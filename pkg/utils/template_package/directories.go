package template_package

import (
	"path/filepath"

	"github.com/rohitramu/kpm/pkg/utils/constants"
	"github.com/rohitramu/kpm/pkg/utils/files"
)

// GetDefaultKpmHomeDir returns the default location of the KPM home directory.
func GetDefaultKpmHomeDir() (string, error) {
	var err error

	var userHomeDir string
	userHomeDir, err = files.GetUserHomeDir()
	if err != nil {
		return "", err
	}

	var result = filepath.Join(userHomeDir, constants.KpmHomeDirName)

	return result, nil
}

// GetPackageRepositoryDir returns the location of the local package repository.
func GetPackageRepositoryDir(kpmHomeDir string) string {
	var packageRepositoryDir = filepath.Join(kpmHomeDir, constants.PackageRepositoryDirName)

	return packageRepositoryDir
}

// GetDefaultOutputDir returns the default path of the root directory for generated files.
func GetDefaultOutputDir(outputParentDir string) string {
	var outputDirPath = filepath.Join(outputParentDir, constants.GeneratedDirName)

	return outputDirPath
}

// GetDefaultOutputName returns the default output name when executing a package.
func GetDefaultOutputName(packageName string, packageVersion string) string {
	return GetPackageFullName(packageName, packageVersion)
}

// GetDefaultExportDir returns the default path of the root directory for exported files.
func GetDefaultExportDir(exportParentDir string) string {
	var outputDirPath = filepath.Join(exportParentDir, constants.ExportDirName)

	return outputDirPath
}

// GetDefaultExportName returns the default name when exporting a package.
func GetDefaultExportName(packageName string, packageVersion string) string {
	return GetPackageFullName(packageName, packageVersion)
}

// GetPackageDir returns the location of a template package in the KPM home directory.
func GetPackageDir(kpmHomeDir string, packageFullName string) string {
	var packageDir = filepath.Join(GetPackageRepositoryDir(kpmHomeDir), packageFullName)

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
