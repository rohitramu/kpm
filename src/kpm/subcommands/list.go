package subcommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"./utils/constants"
	"./utils/files"
	"./utils/logger"
)

// ListCmd lists all packages that are available for use in the given KPM home directory.
func ListCmd(kpmHomeDirArg *string) error {
	var err error

	// Get KPM home directory
	var kpmHomeDir = files.GetAbsolutePathOrDefault(kpmHomeDirArg, files.GetDefaultKpmHomeDir())

	// Get packages directory
	var packageRepositoryDir = filepath.Join(kpmHomeDir, constants.PackageRepositoryDirName)

	// Get all entries in this directory
	var files []os.FileInfo
	files, err = ioutil.ReadDir(packageRepositoryDir)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Print directory names
	for i, file := range files {
		logger.Default.Info.Println(fmt.Sprintf("%d:\t%s", i, file.Name()))
	}

	return nil
}
