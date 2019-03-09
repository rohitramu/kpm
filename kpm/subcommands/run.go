package subcommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"./common"
	"./utils/constants"
	"./utils/files"
	"./utils/log"
	"./utils/templates"
	"./utils/types"
	"./utils/validation"
)

// RunCmd runs the given template package directory and parameters file,
// and then writes the output files to the given output directory.
func RunCmd(packageNameArg *string, packageVersionArg *string, parametersFilePathArg *string, outputNameArg *string, outputDirPathArg *string, kpmHomeDirPathArg *string) error {
	var err error

	// Resolve base paths
	var workingDir string
	workingDir, err = files.GetWorkingDir()
	if err != nil {
		return err
	}

	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePathOrDefaultFunc(kpmHomeDirPathArg, constants.GetDefaultKpmHomeDirPath)
	if err != nil {
		return err
	}

	// Validate name
	var packageName string
	packageName, err = validation.GetStringOrError(packageNameArg, "packageName")
	if err != nil {
		return err
	}

	// Validate version
	var wildcardPackageVersion = validation.GetStringOrDefault(packageVersionArg, "*")

	// Check remote repository for newest matching versions of the package
	var pulledVersion string
	pulledVersion, err = common.PullPackage(packageName, wildcardPackageVersion)
	if err != nil {
		log.Warning(err)
	} else {
		wildcardPackageVersion = pulledVersion
	}

	// Resolve the package version
	var resolvedPackageVersion string
	if resolvedPackageVersion, err = common.ResolvePackageVersion(kpmHomeDir, packageName, wildcardPackageVersion); err != nil {
		return err
	}

	// Resolve generation paths
	var packageFullName = constants.GetPackageFullName(packageName, resolvedPackageVersion)
	var packageDirPath = constants.GetPackageDirPath(constants.GetPackageRepositoryDirPath(kpmHomeDir), packageFullName)
	var outputName = validation.GetStringOrDefault(outputNameArg, packageName)
	var outputParentDir string
	outputParentDir, err = files.GetAbsolutePathOrDefault(outputDirPathArg, workingDir)
	if err != nil {
		return err
	}
	var outputDirPath = constants.GetOutputDirPath(outputParentDir, outputName)
	var parametersFilePath string
	parametersFilePath, err = files.GetAbsolutePathOrDefault(parametersFilePathArg, constants.GetDefaultParametersFilePath(packageDirPath))
	if err != nil {
		return err
	}

	// Log resolved values
	log.Verbose("====")
	log.Verbose(fmt.Sprintf("Package name:      %s", packageName))
	log.Verbose(fmt.Sprintf("Package version:   %s", resolvedPackageVersion))
	log.Verbose(fmt.Sprintf("Package directory: %s", packageDirPath))
	log.Verbose(fmt.Sprintf("Parameters file:   %s", parametersFilePath))
	log.Verbose(fmt.Sprintf("Output name:       %s", outputName))
	log.Verbose(fmt.Sprintf("Output directory:  %s", outputDirPath))
	log.Verbose("====")

	// Get the dependency tree
	var parameters *types.GenericMap
	parameters, err = common.GetPackageParameters(parametersFilePath)
	if err != nil {
		return err
	}
	var dependencyTree *common.DependencyTree
	if dependencyTree, err = common.GetDependencyTree(outputName, kpmHomeDir, packageName, wildcardPackageVersion, parameters); err != nil {
		return err
	}

	// Delete the output directory in case it isn't empty
	os.RemoveAll(outputDirPath)

	// Execute template packages in the dependency tree and write the output to the filesystem
	var numPackages int
	numPackages, err = dependencyTree.VisitNodesDepthFirst(func(pathSegments []string, executableTemplates []*template.Template, templateInput *types.GenericMap) error {
		// Get the output directory
		var outputDir = filepath.Join(outputDirPath, filepath.Join(pathSegments...))

		// Create the output directory if it doesn't exist
		os.MkdirAll(outputDir, os.ModePerm)

		// Get the templates in the package
		for _, tmpl := range executableTemplates {
			// Execute the template with the provided input data
			var templateOutput []byte
			templateOutput, err = templates.ExecuteTemplate(tmpl, templateInput)
			if err != nil {
				return err
			}

			// Write the data to the filesystem
			var outputFilePath = filepath.Join(outputDir, tmpl.Name())
			log.Verbose(fmt.Sprintf("Output file path: %s", outputFilePath))
			ioutil.WriteFile(outputFilePath, templateOutput, os.ModeAppend)
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Verbose(fmt.Sprintf("Executed %d packages", numPackages))

	// Print status
	log.Info(fmt.Sprintf("SUCCESS - Generated output in directory: %s", outputDirPath))

	return nil
}
