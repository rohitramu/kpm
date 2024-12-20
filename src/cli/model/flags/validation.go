package flags

import (
	"fmt"

	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg/utils/files"
)

func CombineValidationFuncs[T any](flagValidationFuncs ...types.FlagIsValidFunc[T]) types.FlagIsValidFunc[T] {
	return func(flagName string, flagValueRef *T) error {
		for _, isValidFunc := range flagValidationFuncs {
			if err := isValidFunc(flagName, flagValueRef); err != nil {
				return err
			}
		}

		// All validations passed.
		return nil
	}
}

func ValidateStringFlagIsSet() types.FlagIsValidFunc[string] {
	return func(flagName string, flagValueRef *string) error {
		if flagValueRef == nil {
			return fmt.Errorf("flag '--%s' must be set", flagName)
		}

		return nil
	}
}

func ValidateDirExists() types.FlagIsValidFunc[string] {
	return func(flagName string, flagValueRef *string) (err error) {
		// Skip this validation if the value isn't set.
		if flagValueRef == nil {
			return nil
		}

		_, err = validateDirExists(flagName, *flagValueRef)
		return err
	}
}

func ValidateDirIsEmpty() types.FlagIsValidFunc[string] {
	return func(flagName string, flagValueRef *string) (err error) {
		// Skip this validation if the value isn't set.
		if flagValueRef == nil {
			return nil
		}

		var absoluteDirPath string
		absoluteDirPath, err = validateDirExists(flagName, *flagValueRef)
		if err != nil {
			return err
		}

		var isEmpty bool
		isEmpty, err = files.DirIsEmpty(absoluteDirPath, flagName)
		if err != nil {
			return err
		}

		if !isEmpty {
			return fmt.Errorf("directory is not empty (flag '--%s')", flagName)
		}

		return nil
	}
}

func validateDirExists(flagName string, dirPath string) (absoluteDirPath string, err error) {
	absoluteDirPath, err = files.GetAbsolutePath(dirPath)
	if err != nil {
		return "", fmt.Errorf("invalid path (flag '--%s'): %s", flagName, err)
	}

	err = files.DirExists(absoluteDirPath, flagName)
	if err != nil {
		return absoluteDirPath, fmt.Errorf("directory doesn't exist (flag '--%s'): %s", flagName, err)
	}

	return absoluteDirPath, nil
}
