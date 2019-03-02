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
	"./utils/logger"
	"./utils/templates"
	"./utils/types"
	"./utils/validation"
)

// GenerateCmd creates Kubernetes configuration files from the
// given template package directory and parameters file, and then
// writes them to the given output directory.
func GenerateCmd(packageNameArg *string, packageVersionArg *string, parametersFilePathArg *string, outputDirPathArg *string, kpmHomeDirPathArg *string) error {
	// Validate string arguments
	var (
		packageName    = validation.GetStringOrFail(packageNameArg, "packageName")
		packageVersion = validation.GetStringOrDefault(packageVersionArg, "*")
	)

	// Resolve base paths
	var (
		workingDir = files.GetWorkingDir()
		kpmHomeDir = files.GetAbsolutePathOrDefault(kpmHomeDirPathArg, files.GetDefaultKpmHomeDir())
	)

	// Check remote repository for newest matching versions of the package
	if pulledVersion, err := common.PullPackage(packageName, packageVersion); err != nil {
		logger.Default.Warning.Println(err)
	} else {
		packageVersion = pulledVersion
	}

	// Resolve generation paths
	var (
		packageDirPath     = common.GetPackageDir(kpmHomeDir, packageName, packageVersion)
		outputDirPath      = files.GetAbsolutePathOrDefault(outputDirPathArg, filepath.Join(workingDir, constants.GeneratedDirName, filepath.Base(packageDirPath)))
		parametersFilePath = files.GetAbsolutePathOrDefault(parametersFilePathArg, filepath.Join(packageDirPath, constants.ParametersFileName))
	)

	// Log resolved paths
	logger.Default.Verbose.Println("====")
	logger.Default.Verbose.Println(fmt.Sprintf("Package name:      %s", packageName))
	logger.Default.Verbose.Println(fmt.Sprintf("Package version:   %s", packageVersion))
	logger.Default.Verbose.Println(fmt.Sprintf("Package directory: %s", packageDirPath))
	logger.Default.Verbose.Println(fmt.Sprintf("Parameters file:   %s", parametersFilePath))
	logger.Default.Verbose.Println(fmt.Sprintf("Output directory:  %s", outputDirPath))
	logger.Default.Verbose.Println("====")

	// Define directory locations
	var (
		dependenciesDirPath = filepath.Join(packageDirPath, constants.DependenciesDirName)
		templatesDirPath    = filepath.Join(packageDirPath, constants.TemplatesDirName)
		helpersDirPath      = filepath.Join(packageDirPath, constants.HelpersDirName)
	)

	// Get template from helpers
	var helpersTemplate, numHelpers = templates.ChainTemplatesFromDir(templates.GetRootTemplate(), helpersDirPath)
	logger.Default.Verbose.Println(fmt.Sprintf("Found %d helper template(s) in directory: %s", numHelpers, helpersDirPath))

	// Get template input values by applying parameters to interface
	var templateInput = common.GetPackageInput(helpersTemplate, packageDirPath, parametersFilePath)

	// Generate output files from dependencies
	processDependenciesAndWriteToFilesystem(dependenciesDirPath, outputDirPath, helpersTemplate, templateInput)

	// Generate output files from the package and write them to the output directory
	var numProcessedTemplates = processTemplatesAndWriteToFilesystem(helpersTemplate, templatesDirPath, templateInput, outputDirPath)
	logger.Default.Verbose.Println(fmt.Sprintf("Processed %d template(s) in directory: %s", numProcessedTemplates, templatesDirPath))

	// Print status
	logger.Default.Info.Println(fmt.Sprintf("SUCCESS - Generated output in directory: %s", outputDirPath))

	return nil
}

// +----------------------+
// | Process dependencies |
// +----------------------+

func processDependenciesAndWriteToFilesystem(dependenciesDirPath string, outputDirPath string, parentTemplate *template.Template, templateInput *types.GenericMap) {

}

// +-------------------+
// | Process templates |
// +-------------------+

func processTemplatesAndWriteToFilesystem(parentTemplate *template.Template, templatesDirPath string, templateInput *types.GenericMap, outputDirPath string) int {
	// Delete and re-create the output directory in case it isn't empty or doesn't exist
	os.RemoveAll(outputDirPath)
	os.MkdirAll(outputDirPath, os.ModePerm)

	// Generate output from each template
	var numTemplates = templates.VisitTemplatesFromDir(templatesDirPath, func() *template.Template {
		// Use the given parent template
		return parentTemplate
	}, func(tmpl *template.Template) {
		// Apply input values to the template
		var generatedFileBytes = templates.ExecuteTemplate(tmpl, templateInput)

		// Write the output to a file
		var outputFilePath = filepath.Join(outputDirPath, tmpl.Name())
		ioutil.WriteFile(outputFilePath, generatedFileBytes, os.ModeAppend|os.ModePerm)
	})

	return numTemplates
}
