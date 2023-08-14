package pkg

import (
	"github.com/rohitramu/kpm/pkg/common"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

// ListCmd lists all packages that are available for use in the given KPM home directory.
func ListCmd(kpmHomeDirPath string) error {
	var err error

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePath(kpmHomeDirPath)
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
		log.Outputf(n)
	}

	return nil
}
