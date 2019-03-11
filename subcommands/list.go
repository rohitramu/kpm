package subcommands

import (
	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
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

	// Get package full names
	var packages []string
	packages, err = common.GetPackageFullNamesFromLocalRepository(kpmHomeDir)
	if err != nil {
		return err
	}

	// Print package full names in order
	for _, n := range packages {
		log.Output(n)
	}

	return nil
}
