package common

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"../utils/constants"
	"../utils/files"
	"../utils/logger"
	"../utils/templates"
	"../utils/types"
	"../utils/validation"
	"../utils/yaml"
)

// PullPackage retrieves a remote template package and makes it available for use.  If a package
// was successfully retrieved, this function returns the retrieved version number.
func PullPackage(packageName string, packageVersion string) (string, error) {
	//TODO: Get list of versions in remote repository

	//TODO: Resolve version to the highest that is compatible with the requested version

	//TODO: Download the template package of the resolved version into the local package repository
	//TODO: Delete the existing package first if it already exists

	return "", fmt.Errorf("Could not find a compatible version for package in remote repository: %s", GetPackageFullName(packageName, packageVersion))
}

// GetTemplateInput creates the input values for a template by combining the interface, parameters and package info.
func GetTemplateInput(parentTemplate *template.Template, packageDirPath string, parameters *types.GenericMap) *types.GenericMap {
	// Add top-level objects
	var result = types.GenericMap{}
	result[constants.TemplateFieldPackage] = GetPackageInfo(packageDirPath)
	result[constants.TemplateFieldValues] = getValuesFromInterface(parentTemplate, packageDirPath, parameters)

	return &result
}

// GetSharedTemplate creates a template which contains default options, functions and
// helper template definitions defined in the given package.
func GetSharedTemplate(packageDirPath string) *template.Template {
	// Get the directory which contains the helper templates
	var helpersDirPath = GetHelpersDirPath(packageDirPath)

	// Create a template which includes the helper template definitions
	var sharedTemplate, numHelpers = templates.ChainTemplatesFromDir(templates.GetRootTemplate(), helpersDirPath)
	logger.Default.Verbose.Println(fmt.Sprintf("Found %d template(s) in directory: %s", numHelpers, helpersDirPath))

	return sharedTemplate
}

// GetPackageInfo returns the package info object for a given package.
func GetPackageInfo(packageDirPath string) *types.PackageInfo {
	var err error

	// Make sure that the package exists
	if fileInfo, err := os.Stat(packageDirPath); err != nil {
		if os.IsNotExist(err) {
			logger.Default.Error.Fatalln(fmt.Sprintf("Package not found in directory: %s", packageDirPath))
		} else {
			logger.Default.Error.Panicln(err)
		}
	} else if !fileInfo.IsDir() {
		logger.Default.Error.Fatalln(fmt.Sprintf("Package path does not point to a directory: %s", packageDirPath))
	} else {
		logger.Default.Verbose.Println(fmt.Sprintf("Found template package in directory: %s", packageDirPath))
	}

	// Check that the package info file exists
	var packageInfoFilePath = filepath.Join(packageDirPath, constants.PackageInfoFileName)
	if fileInfo, err := os.Stat(packageInfoFilePath); err != nil {
		if os.IsNotExist(err) {
			logger.Default.Error.Fatalln(fmt.Sprintf("Package information file does not exist: %s", packageInfoFilePath))
		} else {
			logger.Default.Error.Panicln(err)
		}
	} else if fileInfo.IsDir() {
		logger.Default.Error.Fatalln(fmt.Sprintf("Package information file path does not point to a file: %s", packageInfoFilePath))
	} else {
		logger.Default.Verbose.Println(fmt.Sprintf("Found package information file: %s", packageInfoFilePath))
	}

	// Get package info file content
	var yamlBytes = files.ReadFileToBytes(packageInfoFilePath)

	// Get package info object from file content
	var packageInfo = new(types.PackageInfo)
	yaml.BytesToObject(yamlBytes, packageInfo)

	// Validate package name
	err = validation.ValidatePackageName(packageInfo.Name)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageInfo.Version, false)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return packageInfo
}

// GetPackageParameters returns the parameters in a file as an object which can be used as input to the interface template in a package.
func GetPackageParameters(parametersFilePath string) *types.GenericMap {
	// Get parameters file content as bytes
	var parametersFileBytes []byte
	if fileInfo, err := os.Stat(parametersFilePath); err != nil {
		if os.IsNotExist(err) {
			logger.Default.Warning.Println(fmt.Sprintf("Parameters file does not exist: %s", parametersFilePath))
		} else {
			logger.Default.Error.Panicln(err)
		}
	} else if fileInfo.IsDir() {
		logger.Default.Warning.Println(fmt.Sprintf("Parameters file path does not point to a file: %s", parametersFilePath))
	} else {
		logger.Default.Verbose.Println(fmt.Sprintf("Found parameters file: %s", parametersFilePath))
		parametersFileBytes = files.ReadFileToBytes(parametersFilePath)
	}

	// Get parameters
	var parameters = new(types.GenericMap)
	yaml.BytesToObject(parametersFileBytes, parameters)

	return parameters
}

// GetExecutableTemplates returns all executable templates in a template package.
func GetExecutableTemplates(parentTemplate *template.Template, packageDirPath string) []*template.Template {
	// Get the templates directory
	var executableTemplatesDir = GetTemplatesDirPath(packageDirPath)
	if fileInfo, err := os.Stat(executableTemplatesDir); err != nil {
		if os.IsNotExist(err) {
			logger.Default.Warning.Println(fmt.Sprintf("Template directory does not exist: %s", executableTemplatesDir))
		} else {
			logger.Default.Error.Panicln(err)
		}
	} else if !fileInfo.IsDir() {
		logger.Default.Warning.Println(fmt.Sprintf("Template directory path does not point to a directory: %s", executableTemplatesDir))
	}

	// Return the templates in the directory
	logger.Default.Verbose.Println(fmt.Sprintf("Found template directory: %s", executableTemplatesDir))
	return templates.GetTemplatesFromDir(parentTemplate, executableTemplatesDir)
}

// GetDependencyDefinitionTemplates returns the templates for all dependency definition templates in a template package.
func GetDependencyDefinitionTemplates(parentTemplate *template.Template, packageDirPath string) []*template.Template {
	var dependencyTemplates []*template.Template

	// Get the dependencies directory
	var dependenciesDir = GetDependenciesDirPath(packageDirPath)
	if fileInfo, err := os.Stat(dependenciesDir); err != nil {
		if os.IsNotExist(err) {
			logger.Default.Warning.Println(fmt.Sprintf("Dependency template directory does not exist: %s", dependenciesDir))
		} else {
			logger.Default.Error.Panicln(err)
		}
	} else if !fileInfo.IsDir() {
		logger.Default.Warning.Println(fmt.Sprintf("Dependency template directory path does not point to a directory: %s", dependenciesDir))
	} else {
		logger.Default.Verbose.Println(fmt.Sprintf("Found dependency template directory: %s", dependenciesDir))

		dependencyTemplates = templates.GetTemplatesFromDir(parentTemplate, dependenciesDir)
	}

	return dependencyTemplates
}

// getValuesFromInterface creates the values which can be used as input to templates by executing the interface with parameters.
func getValuesFromInterface(parentTemplate *template.Template, packageDirPath string, parameters *types.GenericMap) *types.GenericMap {
	// Create template object from interface file
	var templateName = constants.InterfaceFileName
	var interfaceFilePath = filepath.Join(packageDirPath, templateName)
	var tmpl = templates.GetTemplateFromFile(parentTemplate, templateName, interfaceFilePath)

	// Generate values by applying parameters to interface
	var interfaceBytes = templates.ExecuteTemplate(tmpl, parameters)

	// Get values object from generated values yaml file
	var result = new(types.GenericMap)
	yaml.BytesToObject(interfaceBytes, result)

	return result
}
