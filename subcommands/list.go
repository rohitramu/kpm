package subcommands

import (
	"strings"

	"./common"
	"./utils/constants"
	"./utils/files"
	"./utils/log"
)

// ListCmd lists all packages that are available for use in the given KPM home directory.
func ListCmd(kpmHomeDirArg *string) error {
	var err error

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirArg, constants.GetDefaultKpmHomeDir)
	if err != nil {
		return err
	}

	var packages []string
	packages, err = common.GetPackageNamesFromLocalRepository(kpmHomeDir)
	if err != nil {
		return err
	}

	// Print directory names
	var output = strings.Join(packages, "\n")
	log.Info("Available template packages:\n" + output)

	return nil
}
