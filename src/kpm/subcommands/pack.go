package subcommands

import (
	"fmt"
	"os"
	"path/filepath"

	"./common"
	"./utils/constants"
	"./utils/files"
	"./utils/logger"
)

// PackCmd packs a local template package so it available for use in the given KPM home directory.
func PackCmd(packageDirPathArg *string, kpmHomeDirPathArg *string) error {
	// Resolve paths
	var (
		packageDirPath                = files.GetAbsolutePathOrDefault(packageDirPathArg, files.GetWorkingDir())
		kpmHomeDirPath                = files.GetAbsolutePathOrDefault(kpmHomeDirPathArg, filepath.Join(files.GetUserHomeDir(), constants.KpmHomeDirName))
		localPackageRepositoryDirPath = filepath.Join(kpmHomeDirPath, constants.PackageRepositoryDirName)
	)

	// Log resolved paths
	logger.Default.Verbose.Println("====")
	logger.Default.Verbose.Println(fmt.Sprintf("Package directory:             %s", packageDirPath))
	logger.Default.Verbose.Println(fmt.Sprintf("Package repository directory:  %s", localPackageRepositoryDirPath))
	logger.Default.Verbose.Println("====")

	// Validate package and get package info
	logger.Default.Verbose.Println("Getting package info")
	var packageInfo = common.GetPackageInfo(packageDirPath)

	// Get package name with version and output path
	var packageNameWithVersion = common.GetPackageFullName(packageInfo.Name, packageInfo.Version)
	var outputDirPath = common.GetPackageDirPath(localPackageRepositoryDirPath, packageNameWithVersion)

	// Delete the output directory in case it isn't empty
	os.RemoveAll(outputDirPath)

	// Copy package to output directory
	logger.Default.Verbose.Println(fmt.Sprintf("Copying package to: %s", outputDirPath))
	files.CopyDir(packageDirPath, outputDirPath)

	logger.Default.Info.Println(fmt.Sprintf("Template package name:    %s", packageInfo.Name))
	logger.Default.Info.Println(fmt.Sprintf("Template package version: %s", packageInfo.Version))

	logger.Default.Info.Println(fmt.Sprintf("SUCCESS - Pack command created package in directory: %s", outputDirPath))

	return nil
}
