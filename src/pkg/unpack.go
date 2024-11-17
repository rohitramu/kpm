package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
)

// UnpackCmd exports a template package to the specified path.
func UnpackCmd(
	packageName string,
	packageVersion string,
	exportDir string,
	exportName string,
	kpmHomeDirPath string,
	userHasConfirmed bool,
) error {
	var err error

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePath(kpmHomeDirPath)
	if err != nil {
		return err
	}

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return err
	}

	// Resolve export paths
	var packageFullName = template_package.GetPackageFullName(packageName, packageVersion)
	var packageDir = template_package.GetPackageDir(kpmHomeDir, packageFullName)

	// Check that the package exists in the local repository
	err = files.DirExists(packageDir, "template package")
	if err != nil {
		return fmt.Errorf("failed to get package \"%s\": %s", packageFullName, err)
	}

	// Log resolved values
	log.Verbosef("====")
	log.Verbosef("Package name:      %s", packageName)
	log.Verbosef("Package version:   %s", packageVersion)
	log.Verbosef("Package directory: %s", packageDir)
	log.Verbosef("Export name:       %s", exportName)
	log.Verbosef("Export directory:  %s", exportDir)
	log.Verbosef("====")

	// Get full export path
	var exportPath = filepath.Join(exportDir, exportName)

	if err = files.DeleteDirIfExists(exportDir, "export", userHasConfirmed); err != nil {
		return err
	}

	err = os.MkdirAll(exportPath, os.ModePerm)
	if err != nil {
		log.Panicf("Failed to create directory: %s\n%s", exportPath, err)
	}

	// Copy package to export path
	log.Debugf("Exporting package contents to: %s", exportPath)
	files.CopyDir(packageDir, exportPath)

	log.Infof(fmt.Sprintf("Package '%s' exported to: %s", packageFullName, exportPath))

	return nil
}
