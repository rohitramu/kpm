package subcommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rohitramu/kpm/subcommands/common"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/docker"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/templates"
	"github.com/rohitramu/kpm/subcommands/utils/types"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
)

// RunCmd runs the given template package directory and parameters file,
// and then writes the output files to the given output directory.
func RunCmd(packageNameArg *string, packageVersionArg *string, parametersFilePathArg *string, outputNameArg *string, outputDirPathArg *string, kpmHomeDirPathArg *string, dockerRegistryArg *string) error {
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

	// Resolve generation paths
	var packageFullName = constants.GetPackageFullName(packageName, packageVersion)
	var packageDirPath = constants.GetPackageDir(kpmHomeDir, packageFullName)
	var outputName = validation.GetStringOrDefault(outputNameArg, constants.GetDefaultOutputName(packageName, packageVersion))
	var outputDirPath string
	outputDirPath, err = files.GetAbsolutePathOrDefault(outputDirPathArg, constants.GetDefaultOutputDir(workingDir))
	if err != nil {
		return err
	}
	var parametersFilePath string
	parametersFilePath, err = files.GetAbsolutePathOrDefault(parametersFilePathArg, constants.GetDefaultParametersFile(packageDirPath))
	if err != nil {
		return err
	}

	// Log resolved values
	log.Info("====")
	log.Info("Package name:      %s", packageName)
	log.Info("Package version:   %s", packageVersion)
	log.Info("Package directory: %s", packageDirPath)
	log.Info("Parameters file:   %s", parametersFilePath)
	log.Info("Output name:       %s", outputName)
	log.Info("Output directory:  %s", outputDirPath)
	log.Info("====")

	// Get the default parameters
	var packageParameters *types.GenericMap
	packageParameters, err = common.GetPackageParameters(parametersFilePath)
	if err != nil {
		return err
	}

	// Get the dependency tree
	var dependencyTree *common.DependencyTree
	if dependencyTree, err = common.GetDependencyTree(kpmHomeDir, packageName, packageVersion, dockerRegistry, outputName, packageParameters); err != nil {
		return err
	}

	// Delete the output directory in case it isn't empty
	err = os.RemoveAll(filepath.Join(outputDirPath, outputName))
	if err != nil {
		return err
	}

	// Execute template packages in the dependency tree and write the output to the filesystem
	var numPackages int
	numPackages, err = dependencyTree.VisitNodesDepthFirst(func(relativeFilePath []string, friendlyNamePath []string, executableTemplates []*template.Template, templateInput *types.GenericMap) error {
		// Get the output directory
		var outputDir = filepath.Join(outputDirPath, filepath.Join(relativeFilePath...))

		// Create the output directory if it doesn't exist
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return err
		}

		// Get the templates in the package
		for _, tmpl := range executableTemplates {
			// Execute the template with the provided input data
			var templateOutput []byte
			templateOutput, err = templates.ExecuteTemplate(tmpl, templateInput)
			if err != nil {
				return fmt.Errorf("Failed to execute package: %s\n%s", strings.Join(friendlyNamePath, " -> "), err)
			}

			// Write the data to the filesystem
			var outputFilePath = filepath.Join(outputDir, tmpl.Name())
			log.Verbose("Output file path: %s", outputFilePath)
			ioutil.WriteFile(outputFilePath, templateOutput, os.ModeAppend)
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Verbose("Executed %d packages", numPackages)

	return nil
}
