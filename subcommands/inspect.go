package subcommands

import (
	"fmt"
	"os"

	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/docker"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// InspectCmd displays the given template package's parameters file.
func InspectCmd(packageNameArg *string, packageVersionArg *string, kpmHomeDirPathArg *string, dockerRegistryArg *string) error {
	var err error

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
			return err
		}
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return err
	}

	// Get Docker registry name
	var dockerRegistry = validation.GetStringOrDefault(dockerRegistryArg, docker.DefaultDockerRegistry)

	// Get parameters file path
	var packageFullName = constants.GetPackageFullName(packageName, packageVersion)
	var packageDirPath = constants.GetPackageDir(kpmHomeDir, packageFullName)
	var parametersFilePath = constants.GetDefaultParametersFile(packageDirPath)

	// Log resolved values
	log.Verbose("====")
	log.Verbose("Package name:      %s", packageName)
	log.Verbose("Package version:   %s", packageVersion)
	log.Verbose("Parameters file:   %s", parametersFilePath)
	log.Verbose("====")

	// Check local repository for package
	var packages []string
	packages, err = common.GetPackageFullNamesFromLocalRepository(kpmHomeDir)
	if err != nil {
		return err
	}
	var found = false
	for _, value := range packages {
		if value == packageFullName {
			found = true
			break
		}
	}

	// If we didn't find the package locally, pull it from the Docker registry
	if !found {
		log.Warning("Package \"%s\" not found in local repository, now checking docker registry \"%s\"...", packageFullName, dockerRegistry)

		// Check remote repository for package
		err = common.PullPackage(kpmHomeDir, dockerRegistry, packageName, packageVersion)
		if err != nil {
			return fmt.Errorf("Failed to get package \"%s\" from docker registry \"%s\": %s", packageFullName, dockerRegistry, err)
		}
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
