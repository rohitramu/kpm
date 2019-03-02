package files

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/otiai10/copy"

	"../constants"
	"../logger"
)

// GetWorkingDir returns the current working directory.
func GetWorkingDir() string {
	var result, err = os.Getwd()
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return result
}

// GetAbsolutePathOrDefault returns the absolute path of the provided path if
// it is not null, otherwise the absolute path of the default path.
func GetAbsolutePathOrDefault(path *string, defaultPath string) string {
	if path == nil {
		return GetAbsolutePath(defaultPath)
	}

	return GetAbsolutePath(*path)
}

// GetAbsolutePath returns the absolute path of the provided path.
func GetAbsolutePath(path string) string {
	var err error

	var outputPath = path

	// Resolve "~" to the user's home directory if required
	if pathSegments := filepath.SplitList(outputPath); len(pathSegments) > 0 && pathSegments[0] == "~" {
		pathSegments[0] = GetUserHomeDir()
		outputPath = filepath.Join(pathSegments...)
	}

	// Check if path is already absolute
	if !filepath.IsAbs(outputPath) {
		// Get absolute path
		outputPath, err = filepath.Abs(outputPath)

		// Exit on error
		if err != nil {
			logger.Default.Error.Fatalln(err)
		}
	}

	return outputPath
}

// GetUserHomeDir returns the path to the home directory of the current user.
func GetUserHomeDir() string {
	var usr, err = user.Current()
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return usr.HomeDir
}

// GetDefaultKpmHomeDir returns the path to the default KPM home directory for the current user.
func GetDefaultKpmHomeDir() string {
	return filepath.Join(GetUserHomeDir(), constants.KpmHomeDirName)
}

// ReadFileToString returns the contents of the given file as a string.
func ReadFileToString(filePath string) string {
	var result = string(ReadFileToBytes(filePath))

	return result
}

// ReadFileToBytes returns the contents of the given file as a byte array.
func ReadFileToBytes(filePath string) []byte {
	var fileData, err = ioutil.ReadFile(filePath)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return fileData
}

// CopyDir recursively copies a directory from source to destination.
func CopyDir(source string, destination string) {
	var err = copy.Copy(source, destination)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}
}
