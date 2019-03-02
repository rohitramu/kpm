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
		logger.Default.Error.Fatalln(fmt.Sprintf("Package path does not point to a file: %s", packageInfoFilePath))
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

// getValuesFromInterface creates the values which can be used as input to templates by executing the interface with parameters.
func getValuesFromInterface(parentTemplate *template.Template, packageDirPath string, parametersFilePath string) *types.GenericMap {
	// Create template object from interface file
	var templateName = constants.InterfaceFileName
	var interfaceFilePath = filepath.Join(packageDirPath, templateName)
	var tmpl = templates.GetTemplateFromFile(parentTemplate, templateName, interfaceFilePath)

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

	// Generate values by applying parameters to interface
	var interfaceBytes = templates.ExecuteTemplate(tmpl, parameters)

	// Get values object from generated values yaml file
	var result = new(types.GenericMap)
	yaml.BytesToObject(interfaceBytes, result)

	return result
}
