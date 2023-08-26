package directories

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/pkg/utils/files"
)

// GetDefaultOutputDir returns the default path of the root directory for generated files.
func GetDefaultOutputDir(outputParentDir string) string {
	var outputDirPath = filepath.Join(outputParentDir, constants.GeneratedDirName)

	return outputDirPath
}

// GetDefaultExportDir returns the default path of the root directory for exported files.
func GetDefaultExportDir(exportParentDir string) string {
	var outputDirPath = filepath.Join(exportParentDir, constants.ExportDirName)

	return outputDirPath
}

func GetKpmHomeDir() (string, error) {
	var err error

	// Try to get the KPM home directory from the environment variable.
	var kpmHomeDir = strings.TrimSpace(os.ExpandEnv("$" + constants.KpmHomeDirEnvVariable))
	if kpmHomeDir != "" {
		var kpmHomeDirAbs string
		kpmHomeDirAbs, err = files.GetAbsolutePath(kpmHomeDir)
		if err != nil {
			return "", fmt.Errorf(
				"invalid directory specified for the \"%s\" environment variable '%s': %s",
				constants.KpmHomeDirEnvVariable,
				kpmHomeDir,
				err,
			)
		}

		kpmHomeDir = kpmHomeDirAbs
	} else {
		// Since the environment variable was empty or not defined, use the default value.
		kpmHomeDir, err = getDefaultKpmHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get default KPM home directory: %s", err)
		}
	}

	return kpmHomeDir, nil
}

// GetDefaultKpmHomeDir returns the default location of the KPM home directory.
func getDefaultKpmHomeDir() (string, error) {
	var err error

	var userHomeDir string
	userHomeDir, err = files.GetUserHomeDir()
	if err != nil {
		return "", err
	}

	var result = filepath.Join(userHomeDir, constants.DefaultKpmHomeDirName)

	return result, nil
}
