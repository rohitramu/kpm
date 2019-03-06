package common

import (
	"fmt"
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
func GetTemplateInput(parentTemplate *template.Template, packageDirPath string, parameters *types.GenericMap) (*types.GenericMap, error) {
	var err error

	// Add top-level objects
	var result = types.GenericMap{}
	result[constants.TemplateFieldPackage], err = GetPackageInfo(packageDirPath)
	if err != nil {
		return nil, err
	}
	result[constants.TemplateFieldValues], err = getValuesFromInterface(parentTemplate, packageDirPath, parameters)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSharedTemplate creates a template which contains default options, functions and
// helper template definitions defined in the given package.
func GetSharedTemplate(packageDirPath string) (*template.Template, error) {
	var err error

	// Get the directory which contains the helper templates
	var helpersDirPath = GetHelpersDirPath(packageDirPath)

	// Create a template which includes the helper template definitions
	var sharedTemplate *template.Template
	var numHelpers int
	sharedTemplate, numHelpers, err = templates.ChainTemplatesFromDir(templates.GetRootTemplate(), helpersDirPath)
	if err != nil {
		return nil, err
	}

	logger.Default.Verbose.Println(fmt.Sprintf("Found %d template(s) in directory: %s", numHelpers, helpersDirPath))

	return sharedTemplate, nil
}

// GetPackageInfo returns the package info object for a given package and validates the package directory.
func GetPackageInfo(packageDirPath string) (*types.PackageInfo, error) {
	var err error

	// Make sure that the package exists
	err = files.DirExists(packageDirPath, "package")
	if err != nil {
		return nil, err
	}

	// Check that the package info file exists
	var packageInfoFilePath = filepath.Join(packageDirPath, constants.PackageInfoFileName)
	err = files.FileExists(packageInfoFilePath, "package information")
	if err != nil {
		return nil, err
	}

	// Get package info file content
	var yamlBytes []byte
	yamlBytes, err = files.ReadBytes(packageInfoFilePath)
	if err != nil {
		return nil, err
	}

	// Get package info object from file content
	var packageInfo = new(types.PackageInfo)
	err = yaml.BytesToObject(yamlBytes, packageInfo)
	if err != nil {
		return nil, err
	}

	// Validate package name
	err = validation.ValidatePackageName(packageInfo.Name)
	if err != nil {
		return nil, err
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageInfo.Version, false)
	if err != nil {
		return nil, err
	}

	// Make sure that the interface file exists
	var interfaceFilePath = filepath.Join(packageDirPath, constants.InterfaceFileName)
	err = files.FileExists(interfaceFilePath, "interface")
	if err != nil {
		return nil, err
	}

	// Make sure that the parameters file exists
	var parametersFilePath = GetDefaultParametersFilePath(packageDirPath)
	err = files.FileExists(parametersFilePath, "parameters")
	if err != nil {
		return nil, err
	}

	return packageInfo, nil
}

// GetPackageParameters returns the parameters in a file as an object which can be used as input to the interface template in a package.
func GetPackageParameters(parametersFilePath string) (*types.GenericMap, error) {
	var err error

	// Make sure that the parameters file exists
	err = files.FileExists(parametersFilePath, "parameters")
	if err != nil {
		return nil, err
	}

	// Get parameters file content as bytes
	var parametersFileBytes []byte
	parametersFileBytes, err = files.ReadBytes(parametersFilePath)
	if err != nil {
		return nil, err
	}

	// Get parameters
	var parameters = new(types.GenericMap)
	err = yaml.BytesToObject(parametersFileBytes, parameters)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}

// GetExecutableTemplates returns all executable templates in a template package.
func GetExecutableTemplates(parentTemplate *template.Template, packageDirPath string) ([]*template.Template, error) {
	var err error

	// Get the templates directory
	var executableTemplatesDir = GetTemplatesDirPath(packageDirPath)
	err = files.FileExists(executableTemplatesDir, "templates")
	if err != nil {
		return nil, err
	}

	// Return the templates in the directory
	logger.Default.Verbose.Println(fmt.Sprintf("Found template directory: %s", executableTemplatesDir))
	var result []*template.Template
	result, err = templates.GetTemplatesFromDir(parentTemplate, executableTemplatesDir)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetDependencyDefinitionTemplates returns the templates for all dependency definition templates in a template package.
func GetDependencyDefinitionTemplates(parentTemplate *template.Template, packageDirPath string) ([]*template.Template, error) {
	var err error

	// Get the dependencies directory
	var dependenciesDir = GetDependenciesDirPath(packageDirPath)
	err = files.DirExists(dependenciesDir, "dependencies")
	if err != nil {
		return nil, err
	}

	var dependencyTemplates []*template.Template
	dependencyTemplates, err = templates.GetTemplatesFromDir(parentTemplate, dependenciesDir)
	if err != nil {
		return nil, err
	}

	return dependencyTemplates, nil
}

// getValuesFromInterface creates the values which can be used as input to templates by executing the interface with parameters.
func getValuesFromInterface(parentTemplate *template.Template, packageDirPath string, parameters *types.GenericMap) (*types.GenericMap, error) {
	var err error

	// Create template object from interface file
	var templateName = constants.InterfaceFileName
	var interfaceFilePath = filepath.Join(packageDirPath, templateName)
	var tmpl *template.Template
	tmpl, err = templates.GetTemplateFromFile(parentTemplate, templateName, interfaceFilePath)
	if err != nil {
		return nil, err
	}

	// Generate values by applying parameters to interface
	var interfaceBytes []byte
	interfaceBytes, err = templates.ExecuteTemplate(tmpl, parameters)
	if err != nil {
		return nil, err
	}

	// Get values object from generated values yaml file
	var result = new(types.GenericMap)
	err = yaml.BytesToObject(interfaceBytes, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
