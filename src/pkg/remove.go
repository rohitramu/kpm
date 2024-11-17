package pkg

import (
	"fmt"

	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
)

// RemoveCmd removes a template package from the local KPM repository.
func RemoveCmd(
	packageName string,
	packageVersion string,
	kpmHomeDirPath string,
	userHasConfirmed bool,
) error {
	var err error

	// Get absolute path of KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePath(kpmHomeDirPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of KPM home directory: %s", kpmHomeDirPath)
	}

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return fmt.Errorf("invalid package name: %s", err)
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return fmt.Errorf("invalid package version: %s", err)
	}

	// Resolve package path
	var packageFullName = template_package.GetPackageFullName(packageName, packageVersion)
	var packageDir = template_package.GetPackageDir(kpmHomeDir, packageFullName)

	// Delete the directory
	if err = files.DeleteDirIfExists(packageDir, "template package", userHasConfirmed); err != nil {
		return err
	}

	return nil
}
