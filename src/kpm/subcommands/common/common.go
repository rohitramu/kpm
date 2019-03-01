package common

import (
	"fmt"
	"io/ioutil"
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

// GetPackageInput creates the input values for a package by combining the interface, parameters and package info.
func GetPackageInput(parentTemplate *template.Template, packageDirPath string, parametersFilePath string) *types.GenericMap {
	// Get package info
	var packageInfo = GetPackageInfo(packageDirPath)

	// Generate the values by populating the interface template with parameters
	var values = getValuesFromInterface(parentTemplate, packageDirPath, parametersFilePath)

	// Add top-level objects
	var result = types.GenericMap{}
	result[constants.TemplateFieldPackage] = packageInfo
	result[constants.TemplateFieldValues] = values

	return &result
}

// GetPackageInfo returns the package info object for a given package.
func GetPackageInfo(packageDirPath string) *types.PackageInfo {
	var err error

	// Make sure that the package exists
	var fileInfo = files.GetFileInfo(packageDirPath)
	if fileInfo == nil {
		logger.Default.Error.Fatalln(fmt.Sprintf("Package not found in directory: %s", packageDirPath))
	} else if !fileInfo.IsDir() {
		logger.Default.Error.Fatalln(fmt.Sprintf("Package path does not point to a directory: %s", packageDirPath))
	}

	// Check that the package info file exists
	var packageInfoFilePath = filepath.Join(packageDirPath, constants.PackageInfoFileName)
	if !files.CheckFileExists(packageInfoFilePath) {
		logger.Default.Error.Fatalln(fmt.Sprintf("Package information file does not exist: %s", packageInfoFilePath))
	}

	// Get package info file content
	var yamlBytes = files.ReadFileToBytes(packageInfoFilePath)

	// Get package info object from file content
	var packageInfo = yaml.BytesToPackageInfo(yamlBytes)

	// Validate the package name and version
	err = validation.ValidatePackageName(packageInfo.Name)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}
	err = validation.ValidatePackageVersion(packageInfo.Version, false)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return packageInfo
}

// getValuesFromInterface creates the values which can be used as input to templates by executing the interface with parameters.
func getValuesFromInterface(parentTemplate *template.Template, packageDirPath string, parametersFilePath string) *types.GenericMap {
	var err error

	// Create template object from interface file
	var templateName = constants.InterfaceFileName
	var interfaceFilePath = filepath.Join(packageDirPath, templateName)
	var tmpl = templates.GetTemplateFromFile(parentTemplate, templateName, interfaceFilePath)

	// Get parameters file content as bytes
	var parametersFileBytes []byte
	parametersFileBytes, err = ioutil.ReadFile(parametersFilePath)
	if err != nil {
		logger.Default.Warning.Println(fmt.Sprintf("Failed to read parameters file: %s", err))
		parametersFileBytes = []byte{}
	}

	// Get parameters
	var parameters = yaml.BytesToMap(parametersFileBytes)

	// Generate values by applying parameters to interface
	var interfaceBytes = templates.ExecuteTemplate(tmpl, parameters)

	// Get values object from generated values yaml file
	var result = yaml.BytesToMap(interfaceBytes)

	return result
}
