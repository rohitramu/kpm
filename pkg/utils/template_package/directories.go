package template_package

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rohitramu/kpm/pkg/utils/constants"
	"github.com/rohitramu/kpm/pkg/utils/files"
)

const KpmHomeDirEnvVar = "KPM_HOME"

// GetDefaultKpmHomeDir returns the default location of the KPM home directory.
func getDefaultKpmHomeDir() (string, error) {
	var err error

	var userHomeDir string
	userHomeDir, err = files.GetUserHomeDir()
	if err != nil {
		return "", err
	}

	var result = filepath.Join(userHomeDir, constants.KpmHomeDirName)

	return result, nil
}

func GetKpmHomeDir() (string, error) {
	var err error

	// Try to get the KPM home directory from the environment variable.
	var kpmHomeDir = strings.TrimSpace(os.ExpandEnv("$" + KpmHomeDirEnvVar))
	if kpmHomeDir != "" {
		var kpmHomeDirAbs string
		kpmHomeDirAbs, err = files.GetAbsolutePath(kpmHomeDir)
		if err != nil {
			return "", fmt.Errorf(
				"invalid directory specified for the \"%s\" environment variable '%s': %s",
				KpmHomeDirEnvVar,
				kpmHomeDir,
				err,
			)
		}

		kpmHomeDir = kpmHomeDirAbs
	} else {
		// Since the environment variable was empty or not defined, use the default value.
		kpmHomeDir, err = getDefaultKpmHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get default KPM home directory: %s", err)
		}
	}

	return kpmHomeDir, nil
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
