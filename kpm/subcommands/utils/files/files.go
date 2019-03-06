package files

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

// GetWorkingDir returns the current working directory.
func GetWorkingDir() (string, error) {
	var result, err = os.Getwd()
	if err != nil {
		return "", err
	}

	return result, nil
}

// GetAbsolutePathOrDefault returns the absolute path of the provided path if
// it is not null, otherwise the absolute path of the default path.
func GetAbsolutePathOrDefault(path *string, defaultPath string) (string, error) {
	return GetAbsolutePathOrDefaultFunc(path, func() (string, error) {
		return defaultPath, nil
	})
}

// GetAbsolutePathOrDefaultFunc returns the absolute path of the provided path if
// it is not null, otherwise the absolute path of the default path which is supplied
// by the default path function.
func GetAbsolutePathOrDefaultFunc(path *string, defaultPathFunc func() (string, error)) (string, error) {
	var err error

	// Check if we can just use the provided path, otherwise we need to use the default
	if path != nil {
		// Get absolute path
		var result string
		result, err = GetAbsolutePath(*path)
		if err != nil {
			return "", err
		}

		return result, nil
	}

	// Run provided supplier function
	var defaultPath string
	defaultPath, err = defaultPathFunc()
	if err != nil {
		return "", err
	}

	// Get absolute path of default value
	var result string
	result, err = GetAbsolutePath(defaultPath)
	if err != nil {
		return "", err
	}

	return result, nil
}

// GetAbsolutePath returns the absolute path of the provided path.
func GetAbsolutePath(path string) (string, error) {
	var err error

	var outputPath = path

	// Resolve "~" to the user's home directory if required
	if strings.HasPrefix(outputPath, "~") {
		var usrHomeDir string
		usrHomeDir, err = GetUserHomeDir()
		if err != nil {
			return "", err
		}
		outputPath = usrHomeDir + strings.TrimPrefix(outputPath, "~")
	}

	// Check if path is already absolute
	if !filepath.IsAbs(outputPath) {
		// Get absolute path
		outputPath, err = filepath.Abs(outputPath)

		// Exit on error
		if err != nil {
			return "", err
		}
	}

	return outputPath, nil
}

// GetUserHomeDir returns the path to the home directory of the current user.
func GetUserHomeDir() (string, error) {
	var usr, err = user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

// ReadString returns the contents of the given file as a string.
func ReadString(filePath string) (string, error) {
	var resultBytes, err = ReadBytes(filePath)
	if err != nil {
		return "", err
	}

	var resultString = string(resultBytes)

	return resultString, nil
}

// ReadBytes returns the contents of the given file as a byte array.
func ReadBytes(filePath string) ([]byte, error) {
	var fileData, err = ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

// CopyDir recursively copies a directory from source to destination.
func CopyDir(source string, destination string) error {
	var err = copy.Copy(source, destination)
	if err != nil {
		return err
	}

	return nil
}
