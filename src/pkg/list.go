package pkg

import (
	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
)

// ListCmd lists all packages that are available for use in the given KPM home directory.
func ListCmd(kpmHomeDirPath string) ([]string, error) {
	var err error

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePath(kpmHomeDirPath)
	if err != nil {
		return nil, err
	}

	// Get package full names
	var packages []string
	packages, err = template_package.GetPackageFullNamesFromLocalRepository(kpmHomeDir)
	if err != nil {
		return nil, err
	}

	return packages, nil
}
