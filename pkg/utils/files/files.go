package files

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"

	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/user_prompts"
)

// GetWorkingDir returns the current working directory.
func GetWorkingDir() (string, error) {
	var result, err = os.Getwd()
	if err != nil {
		return "", err
	}

	return result, nil
}

// GetUserHomeDir returns the path to the home directory of the current user.
func GetUserHomeDir() (string, error) {
	var usr, err = user.Current()
	if err != nil {
		return "", err
	}

	var userHomeDir string
	userHomeDir, err = GetAbsolutePath(usr.HomeDir)
	if err != nil {
		return "", nil
	}

	return userHomeDir, nil
}

// GetTempDir returns the path to the system's temporary directory.
func GetTempDir() (string, error) {
	return os.TempDir(), nil
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
	var result = path

	// Resolve "~" to the user's home directory if required
	if IsAbsFromHome(result) {
		result, err = GetExpandedHomeDirPath(result)
		if err != nil {
			return "", err
		}
	}

	// Check if path is already absolute
	if !filepath.IsAbs(result) {
		// Get absolute path
		result, err = filepath.Abs(result)

		// Exit on error
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func IsAbsFromHomeOrRoot(path string) bool {
	return IsAbsFromHome(path) || filepath.IsAbs(path)
}

func IsAbsFromHome(path string) bool {
	if path == "~" || strings.HasPrefix(path, "~/") {
		return true
	}

	return false
}

func GetExpandedHomeDirPath(path string) (string, error) {
	if !IsAbsFromHome(path) {
		return path, nil
	}

	if path == "~" {
		return GetUserHomeDir()
	}

	if strings.HasPrefix(path, "~/") {
		homeDir, err := GetUserHomeDir()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s/%s", homeDir, path[2:]), nil
	}

	// In all other cases, return the path untouched.
	return path, nil
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
	var fileData, err = os.ReadFile(filePath)
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

// FileExists checks whether a file exists and returns an error if it doesn't.
func FileExists(absoluteFilePath string, lowercaseHumanFriendlyName string) error {
	if fileInfo, err := os.Stat(absoluteFilePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s file does not exist: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteFilePath)
		}

		// File may exist, but we had an unexpected failure
		log.Panicf("Failed to search for file: %s", err)
	} else if fileInfo.IsDir() {
		return fmt.Errorf("%s file path does not point to a file: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteFilePath)
	}

	log.Debugf("Found %s file: %s", lowercaseHumanFriendlyName, absoluteFilePath)

	return nil
}

// CreateFile creates a new file.
func CreateFile(absoluteFilePath string, lowercaseHumanFriendlyName string, content string) error {
	var err error

	// Make sure the file doesn't already exist.
	if err = FileExists(absoluteFilePath, lowercaseHumanFriendlyName); err == nil {
		return fmt.Errorf("failed to create %s file - file already exists: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteFilePath)
	}

	// Create any directories in the path that don't exist.
	if err = os.MkdirAll(path.Dir(absoluteFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create parent directory of %s file: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteFilePath)
	}

	// Create the file.
	var f *os.File
	if f, err = os.Create(absoluteFilePath); err != nil {
		return fmt.Errorf("failed to create %s file: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteFilePath)
	}
	defer f.Close()

	// Write content to the file.
	if _, err = f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write content to %s file: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteFilePath)
	}

	return nil
}

// DirExists checks whether a directory exists and returns an error if it doesn't.
func DirExists(absoluteDirPath string, lowercaseHumanFriendlyName string) error {
	if fileInfo, err := os.Stat(absoluteDirPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s directory does not exist: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteDirPath)
		}

		// Directory may exist, but we had an unexpected failure
		log.Panicf("Failed to search for directory: %s", err)
	} else if !fileInfo.IsDir() {
		return fmt.Errorf("%s directory path does not point to a directory: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteDirPath)
	}

	log.Debugf("Found %s directory: %s", lowercaseHumanFriendlyName, absoluteDirPath)

	return nil
}

func DirIsEmpty(absoluteDirPath string, lowercaseHumanFriendlyName string) (bool, error) {
	var err error

	if err = DirExists(absoluteDirPath, lowercaseHumanFriendlyName); err != nil {
		return false, err
	}

	var fileInfo *os.File
	if fileInfo, err = os.Open(absoluteDirPath); err != nil {
		return false, err
	}
	defer fileInfo.Close()

	var names []string
	if names, err = fileInfo.Readdirnames(1); err != io.EOF {
		return false, err
	}

	return len(names) == 0, nil
}

func DeleteDirIfExists(absoluteDirPath string, lowercaseHumanFriendlyName string, userHasConfirmed bool) error {
	var err error

	// If the directory doesn't exist, exit
	if err = DirExists(absoluteDirPath, lowercaseHumanFriendlyName); err != nil {
		return nil
	}

	var dirIsEmpty bool
	if dirIsEmpty, err = DirIsEmpty(absoluteDirPath, lowercaseHumanFriendlyName); err != nil {
		return err
	}

	// If the directory is not empty, get user confirmation
	if !dirIsEmpty {
		// If the user hasn't already confirmed, ask for a confirmation now
		if !userHasConfirmed {
			if userHasConfirmed, err = user_prompts.ConfirmWithUser("%s directory exists, so it will be deleted before continuing: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteDirPath); err != nil {
				return err
			}
		}

		// Couldn't get user confirmation, so return an error saying that the operation has been cancelled
		if !userHasConfirmed {
			return fmt.Errorf("operation cancelled - user did not confirm deletion of pre-existing %s folder", lowercaseHumanFriendlyName)
		}
	}

	// Delete the directory
	if err = os.RemoveAll(absoluteDirPath); err != nil {
		log.Panicf("Failed to delete directory: %s\n%s", absoluteDirPath, err)
	}

	return nil
}

func CreateDir(absoluteDirPath string, lowercaseHumanFriendlyName string, userHasConfirmed bool) error {
	var err error

	// If the directory exists already, exit.
	if err = DirExists(absoluteDirPath, lowercaseHumanFriendlyName); err == nil {
		return fmt.Errorf("directory \"%s\" already exists", absoluteDirPath)
	}

	// Get confirmation from the user if needed.
	if !userHasConfirmed {
		if userHasConfirmed, err = user_prompts.ConfirmWithUser("%s directory will be created before continuing: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteDirPath); err != nil {
			return err
		}
	}

	// Check if user declined confirmation.
	if !userHasConfirmed {
		return fmt.Errorf("operation cancelled - user did not confirm creation of %s directory: %s", toTitleCase(lowercaseHumanFriendlyName), absoluteDirPath)
	}

	// Create the directory.
	if err = os.MkdirAll(absoluteDirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create %s directory: %s", toTitleCase(lowercaseHumanFriendlyName), err)
	}

	return nil
}

func CreateDirIfNotExists(absoluteDirPath string, lowercaseHumanFriendlyName string, userHasConfirmed bool) (err error) {
	err = DirExists(absoluteDirPath, lowercaseHumanFriendlyName)
	if err != nil {
		err = CreateDir(absoluteDirPath, lowercaseHumanFriendlyName, userHasConfirmed)
		if err != nil {
			return err
		}
	}

	return nil
}

func toTitleCase(text string) string {
	return strings.ToUpper(string(text[0])) + text[1:]
}
