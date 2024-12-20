package directories

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
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

func GetOrCreateKpmHomeDir(userHasConfirmed bool) (kpmHomeDir string, err error) {
	kpmHomeDir, err = GetKpmHomeDir()
	if err != nil {
		return "", err
	}

	err = template_package.CreateFilesystemRepoDir(kpmHomeDir, "KPM home", userHasConfirmed)
	if err != nil {
		return "", err
	}

	return kpmHomeDir, nil
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
