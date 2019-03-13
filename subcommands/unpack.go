package subcommands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/docker"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// UnpackCmd exports a template package to the specified path.
func UnpackCmd(packageNameArg *string, packageVersionArg *string, exportDirPathArg *string, exportNameArg *string, kpmHomeDirPathArg *string, dockerRegistryArg *string) error {
	var err error

	// Resolve base paths
	var workingDir string
	workingDir, err = files.GetWorkingDir()
	if err != nil {
		return err
	}

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDir)
	if err != nil {
		return err
	}

	// Get package name
	var packageName string
	packageName, err = validation.GetStringOrError(packageNameArg, "packageName")
	if err != nil {
		return err
	}

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Get package version
	var packageVersion string
	packageVersion, err = validation.GetStringOrError(packageVersionArg, "packageVersion")
	if err != nil {
		// Since the package version was not provided, check the local repository for the highest version
		if packageVersion, err = common.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
			return fmt.Errorf("Package version must be provided if the package does not exist in the local repository: %s", err)
		}
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return err
	}

	// Get Docker registry name
	var dockerRegistry = validation.GetStringOrDefault(dockerRegistryArg, docker.DefaultDockerRegistry)

	// Resolve export paths
	var packageFullName = constants.GetPackageFullName(packageName, packageVersion)
	var packageDir = constants.GetPackageDir(kpmHomeDir, packageFullName)
	var exportName = validation.GetStringOrDefault(exportNameArg, constants.GetDefaultExportName(packageName, packageVersion))
	var exportDir string
	exportDir, err = files.GetAbsolutePathOrDefault(exportDirPathArg, constants.GetDefaultExportDir(workingDir))
	if err != nil {
		return err
	}

	// Check that the package exists in the local repository
	err = files.DirExists(packageDir, "template package")
	if err != nil {
		log.Warning("Package \"%s\" not found in local repository, now checking docker registry \"%s\"...", packageFullName, dockerRegistry)

		// Check remote repository for package
		err = common.PullPackage(kpmHomeDir, dockerRegistry, packageName, packageVersion)
		if err != nil {
			return fmt.Errorf("Failed to get package \"%s\" from docker registry \"%s\": %s", packageFullName, dockerRegistry, err)
		}
	}

	// Log resolved values
	log.Info("====")
	log.Info("Package name:      %s", packageName)
	log.Info("Package version:   %s", packageVersion)
	log.Info("Package directory: %s", packageDir)
	log.Info("Export name:       %s", exportName)
	log.Info("Export directory:  %s", exportDir)
	log.Info("====")

	// Get full export path
	var exportPath = filepath.Join(exportDir, exportName)

	// Delete the export path in case it isn't empty
	err = os.RemoveAll(exportPath)
	if err != nil {
		log.Panic("Failed to remove directory: %s\n%s", exportPath, err)
	}

	err = os.MkdirAll(exportPath, os.ModePerm)
	if err != nil {
		log.Panic("Failed to create directory: %s\n%s", exportPath, err)
	}

	// Copy package to export path
	log.Verbose("Exporting package to: %s", exportPath)
	files.CopyDir(packageDir, exportPath)

	return nil
}
