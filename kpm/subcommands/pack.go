package subcommands

import (
	"os"

	"./common"
	"./utils/constants"
	"./utils/files"
	"./utils/log"
	"./utils/types"
)

// PackCmd packs a local template package so it is available for use in the given local KPM repository.
func PackCmd(packageDirPathArg *string, kpmHomeDirPathArg *string) error {
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

	// Get local package repository directory
	var localPackageRepositoryDir = constants.GetPackageRepositoryDir(kpmHomeDir)

	// Log resolved paths
	log.Info("====")
	log.Info("Package directory:             %s", packageDir)
	log.Info("Package repository directory:  %s", localPackageRepositoryDir)
	log.Info("====")

	// Validate package and get package info
	log.Verbose("Getting package info")
	var packageInfo *types.PackageInfo
	packageInfo, err = common.GetPackageInfo(kpmHomeDir, packageDir)
	if err != nil {
		return err
	}

	// Get package name with version and output path
	var packageNameWithVersion = constants.GetPackageFullName(packageInfo.Name, packageInfo.Version)
	var outputDir = constants.GetPackageDir(kpmHomeDir, packageNameWithVersion)

	// Delete the output directory in case it isn't empty
	os.RemoveAll(outputDir)

	// Copy package to output directory
	log.Verbose("Copying package to: %s", outputDir)
	files.CopyDir(packageDir, outputDir)

	log.Info("Template package name:    %s", packageInfo.Name)
	log.Info("Template package version: %s", packageInfo.Version)

	log.Verbose("Repository directory:     %s", outputDir)

	return nil
}
