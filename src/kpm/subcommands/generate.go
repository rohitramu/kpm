package subcommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"./common"
	"./utils/files"
	"./utils/logger"
	"./utils/templates"
	"./utils/types"
	"./utils/validation"
)

// GenerateCmd creates Kubernetes configuration files from the
// given template package directory and parameters file, and then
// writes them to the given output directory.
func GenerateCmd(packageNameArg *string, packageVersionArg *string, parametersFilePathArg *string, outputNameArg *string, outputDirPathArg *string, kpmHomeDirPathArg *string) error {
	var err error

	// Validate string arguments
	var (
		packageName            = validation.GetStringOrFail(packageNameArg, "packageName")
		wildcardPackageVersion = validation.GetStringOrDefault(packageVersionArg, "*")
	)

	// Resolve base paths
	var (
		workingDir = files.GetWorkingDir()
		kpmHomeDir = files.GetAbsolutePathOrDefault(kpmHomeDirPathArg, files.GetDefaultKpmHomeDir())
	)

	// Check remote repository for newest matching versions of the package
	var pulledVersion string
	pulledVersion, err = common.PullPackage(packageName, wildcardPackageVersion)
	if err != nil {
		logger.Default.Warning.Println(err)
	} else {
		wildcardPackageVersion = pulledVersion
	}

	// Resolve the package version
	var resolvedPackageVersion string
	if resolvedPackageVersion, err = common.ResolvePackageVersion(kpmHomeDir, packageName, wildcardPackageVersion); err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Resolve generation paths
	var (
		packageFullName    = common.GetPackageFullName(packageName, resolvedPackageVersion)
		packageDirPath     = common.GetPackageDirPath(common.GetPackageRepositoryDirPath(kpmHomeDir), packageFullName)
		outputName         = validation.GetStringOrDefault(outputNameArg, packageFullName)
		outputDirPath      = common.GetOutputDirPath(files.GetAbsolutePathOrDefault(outputDirPathArg, workingDir), outputName)
		parametersFilePath = files.GetAbsolutePathOrDefault(parametersFilePathArg, common.GetDefaultParametersFilePath(packageDirPath))
	)

	// Log resolved values
	logger.Default.Verbose.Println("====")
	logger.Default.Verbose.Println(fmt.Sprintf("Package name:      %s", packageName))
	logger.Default.Verbose.Println(fmt.Sprintf("Package version:   %s", resolvedPackageVersion))
	logger.Default.Verbose.Println(fmt.Sprintf("Package directory: %s", packageDirPath))
	logger.Default.Verbose.Println(fmt.Sprintf("Parameters file:   %s", parametersFilePath))
	logger.Default.Verbose.Println(fmt.Sprintf("Output name:       %s", outputName))
	logger.Default.Verbose.Println(fmt.Sprintf("Output directory:  %s", outputDirPath))
	logger.Default.Verbose.Println("====")

	// Get the dependency tree
	var parameters = common.GetPackageParameters(parametersFilePath)
	var dependencyTree *common.DependencyTree
	if dependencyTree, err = common.GetDependencyTree(outputName, kpmHomeDir, packageName, wildcardPackageVersion, parameters); err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Delete the output directory in case it isn't empty
	os.RemoveAll(outputDirPath)

	// Execute template packages in the dependency tree and write the output to the filesystem
	dependencyTree.VisitNodesDepthFirst(func(pathSegments []string, executableTemplates []*template.Template, templateInput *types.GenericMap) {
		// Get the output directory
		var outputDir = filepath.Join(outputDirPath, filepath.Join(pathSegments...))

		// Create the output directory if it doesn't exist
		os.MkdirAll(outputDir, os.ModePerm)

		// Get the templates in the package
		for _, tmpl := range executableTemplates {
			// Execute the template with the provided input data
			var templateOutput = templates.ExecuteTemplate(tmpl, templateInput)

			// Write the data to the filesystem
			var outputFilePath = filepath.Join(outputDir, tmpl.Name())
			logger.Default.Verbose.Println(fmt.Sprintf("Output file path: %s", outputFilePath))
			ioutil.WriteFile(outputFilePath, templateOutput, os.ModeAppend)
		}
	})

	// Print status
	logger.Default.Info.Println(fmt.Sprintf("SUCCESS - Generated output in directory: %s", outputDirPath))

	return nil
}
