package pkg

import (
	"fmt"
	"path/filepath"

	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
)

func NewTemplatePackageCmd(
	packageName string,
	packagePath string,
	userHasConfirmed bool,
) error {
	var err error

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Get package directory
	var packageDirAbsolutePath string
	packageDirAbsolutePath, err = files.GetAbsolutePath(packagePath)
	if err != nil {
		return err
	}

	// Get full package path
	var packageAbsolutePath = filepath.Join(packageDirAbsolutePath, packageName)

	// Create the template package directory
	if err = files.CreateDir(packageAbsolutePath, "new template package", userHasConfirmed); err != nil {
		return fmt.Errorf("failed to create template package directory: %s", err)
	}

	// Create the template package directory structure
	if err = template_package.GenerateSampleTemplatePackage(packageAbsolutePath, packageName); err != nil {
		return fmt.Errorf("failed to generate template package directory structure: %s", err)
	}

	return nil
}
