package utils

import (
	"os"
	"os/user"
	"path/filepath"
)

var logger = NewLogger()

// GetCurrentWorkingDir returns the current working directory
func GetCurrentWorkingDir() string {
	var result, err = os.Getwd()
	if err != nil {
		logger.Error.Fatalln(err)
	}

	return result
}

// GetAbsolutePathOrDefault returns the absolute path of the provided path if
// it is not null, otherwise the absolute path of the default path
func GetAbsolutePathOrDefault(path *string, defaultPath string) string {
	if path == nil {
		return GetAbsolutePath(defaultPath)
	}

	return GetAbsolutePath(*path)
}

// GetAbsolutePath returns the absolute path of the provided path
func GetAbsolutePath(path string) string {
	var err error

	var outputPath = path

	// Resolve "~" to the user's home directory if required
	if len(outputPath) > 0 && outputPath[0] == '~' && (outputPath[1] == '/' || outputPath[1] == '\\') {
		var usr *(user.User)
		usr, err = user.Current()
		outputPath = filepath.Join(usr.HomeDir, outputPath[2:])
	}

	// Check if path is already absolute
	if !filepath.IsAbs(outputPath) {
		// Get absolute path
		outputPath, err = filepath.Abs(outputPath)

		// Exit on error
		if err != nil {
			logger.Error.Fatalln(err)
		}
	}

	return outputPath
}
