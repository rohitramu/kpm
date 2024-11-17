package pkg

import (
	"fmt"
	"os"

	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
	"golang.org/x/exp/slices"
)

// InspectCmd displays the given template package's parameters file.
func InspectCmd(
	packageName string,
	packageVersion string,
	kpmHomeDirPath string,
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

	// Get parameters file path
	var packageFullName = template_package.GetPackageFullName(packageName, packageVersion)
	var packageDirPath = template_package.GetPackageDir(kpmHomeDir, packageFullName)
	var parametersFilePath = template_package.GetDefaultParametersFile(packageDirPath)

	// Log resolved values
	log.Verbosef("====")
	log.Verbosef("Package name:      %s", packageName)
	log.Verbosef("Package version:   %s", packageVersion)
	log.Verbosef("Parameters file:   %s", parametersFilePath)
	log.Verbosef("====")

	// Check local repository for package
	var packages []string
	packages, err = template_package.GetPackageFullNamesFromLocalRepository(kpmHomeDir)
	if err != nil {
		return err
	}
	if !slices.Contains(packages, packageFullName) {
		return fmt.Errorf("failed to get package \"%s\": %s", packageFullName, err)
	}

	// Get the contents of the default parameters file
	var parametersFile *os.File
	parametersFile, err = os.Open(parametersFilePath)
	if err != nil {
		return err
	}

	// Make sure to close the file afterwards
	defer parametersFile.Close()

	// Print the contents of the default parameters file to output
	log.OutputStream(parametersFile)

	return nil
}
