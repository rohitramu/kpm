package subcommands

import (
	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/types"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// PackCmd packs a local template package so it is available for use in the given local KPM repository.
func PackCmd(packageDirPathArg *string, kpmHomeDirPathArg *string, userHasConfirmedArg *bool) error {
	var err error

	// Package directory
	var packageDir string
	packageDir, err = files.GetAbsolutePathOrDefaultFunc(packageDirPathArg, files.GetWorkingDir)
	if err != nil {
		return err
	}

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDir)
	if err != nil {
		return err
	}

	// Log resolved paths
	log.Info("====")
	log.Info("Package directory:             %s", packageDir)
	log.Info("====")

	// Validate package and get package info
	log.Debug("Getting package info")
	var packageInfo *types.PackageInfo
	packageInfo, err = common.GetPackageInfo(kpmHomeDir, packageDir)
	if err != nil {
		return err
	}

	// Get package name with version and output path
	var packageNameWithVersion = constants.GetPackageFullName(packageInfo.Name, packageInfo.Version)
	var outputDir = constants.GetPackageDir(kpmHomeDir, packageNameWithVersion)

	// Delete the output directory
	var userHasConfirmed bool = validation.GetBoolOrDefault(userHasConfirmedArg, false)
	if err = files.DeleteDirIfExists(outputDir, "output", userHasConfirmed); err != nil {
		return err
	}

	// Copy package to output directory
	log.Debug("Copying package to: %s", outputDir)
	files.CopyDir(packageDir, outputDir)

	log.Info("====")
	log.Info("Template package name:    %s", packageInfo.Name)
	log.Info("Template package version: %s", packageInfo.Version)
	log.Info("Output directory:         %s", outputDir)
	log.Info("====")

	return nil
}
