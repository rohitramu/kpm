package pkg

import (
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

// PackCmd packs a local template package so it is available for use in the given local KPM repository.
func PackCmd(
	packageDirPath string,
	kpmHomeDirPath string,
	userHasConfirmed bool,
) error {
	var err error

	// Package directory
	var packageDir string
	packageDir, err = files.GetAbsolutePath(packageDirPath)
	if err != nil {
		return err
	}

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePath(kpmHomeDirPath)
	if err != nil {
		return err
	}

	// Log resolved paths
	log.Verbosef("====")
	log.Verbosef("Template package directory:             %s", packageDir)
	log.Verbosef("====")

	// Validate package and get package info
	log.Debugf("Getting template package info")
	var packageInfo *templates.PackageInfo
	packageInfo, err = template_package.GetPackageInfo(kpmHomeDir, packageDir)
	if err != nil {
		return err
	}

	// Get package name with version and output path
	var packageNameWithVersion = template_package.GetPackageFullName(packageInfo.Name, packageInfo.Version)
	var outputDir = template_package.GetPackageDir(kpmHomeDir, packageNameWithVersion)

	// Delete the package directory
	if err = files.DeleteDirIfExists(outputDir, "template package", userHasConfirmed); err != nil {
		return err
	}

	// Copy package to output directory
	log.Debugf("Copying package to: %s", outputDir)
	files.CopyDir(packageDir, outputDir)

	log.Verbosef("====")
	log.Verbosef("Template package name:    %s", packageInfo.Name)
	log.Verbosef("Template package version: %s", packageInfo.Version)
	log.Verbosef("Output directory:         %s", outputDir)
	log.Verbosef("====")

	return nil
}
