package subcommands

import (
	"path/filepath"
	"strings"

	"./common"
	"./utils/constants"
	"./utils/files"
	"./utils/logger"
)

// ListCmd lists all packages that are available for use in the given KPM home directory.
func ListCmd(kpmHomeDirArg *string) error {
	var err error

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirArg, constants.GetDefaultKpmHomeDirPath)
	if err != nil {
		return err
	}

	// Get the packages directory
	var packageRepositoryDir = filepath.Join(kpmHomeDir, constants.PackageRepositoryDirName)

	var packages []string
	packages, err = common.GetPackageNamesFromLocalRepository(packageRepositoryDir)
	if err != nil {
		return err
	}

	// Print directory names
	var output = strings.Join(packages, "\n")
	logger.Default.Info.Println("Packages:\n" + output)

	return nil
}
