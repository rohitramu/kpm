package subcommands

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Logger
var logger = log.New(os.Stderr, "", log.LstdFlags)

// GenerateCmd generates a Kubernetes configuration from the given
// template package directory, parameters file and output directory
func GenerateCmd(packageDirPath *string, parametersFilePath *string, outputDirPath *string) error {
	logger.Println(fmt.Sprintf("Package directory: %s", *packageDirPath))
	logger.Println(fmt.Sprintf("Parameters file: %s", *parametersFilePath))
	logger.Println(fmt.Sprintf("Output directory: %s", *outputDirPath))

	// Define well-known file names
	//var dependenciesFileName = "dependencies.yaml"
	var interfaceFileName = "interface.yaml"
	//var packageFileName = "package.yaml"

	// Define file locations
	//var dependenciesFilePath = filepath.Join(*packageDirPath, dependenciesFileName)
	var interfaceFilePath = filepath.Join(*packageDirPath, interfaceFileName)
	//var packageFilePath = filepath.Join(*packageDirPath, packageFileName)

	// Define well-known directory names
	//var dependenciesDirName = "dependencies"
	var templatesDirName = "templates"
	var helpersDirName = "helpers"

	// Define directory locations
	//var dependenciesDirPath = filepath.Join(packageDirPath, dependenciesDirName)
	var templatesDirPath = filepath.Join(*packageDirPath, templatesDirName)
	var helpersDirPath = filepath.Join(*packageDirPath, helpersDirName)

	// Process helpers
	var helpersTemplate = getHelpersTemplate(&helpersDirPath)

	// Generate the interface file from the template and user parameters
	var parameters = getParametersFromInterface(&interfaceFileName, &interfaceFilePath, parametersFilePath)

	processTemplatesAndWriteToFilesystem(helpersTemplate, &templatesDirPath, outputDirPath, parameters)

	logger.Println("Success!")

	return nil
}

func getParametersFromInterface(interfaceFileName *string, interfaceFilePath *string, userParametersFilePath *string) *map[interface{}]interface{} {
	// Get interface template as a string
	var templateString = string(*readFileToBytes(interfaceFilePath))

	// Get user parameters
	var userParameters = yamlBytesToObject(readFileToBytes(userParametersFilePath))

	// Generate interface from template and user parameters
	var interfaceBytes = executeTemplate(nil, interfaceFileName, &templateString, userParameters)

	// Get values from generated interface
	var result = yamlBytesToObject(interfaceBytes)

	return result
}

func getHelpersTemplate(helpersDirPath *string) *template.Template {
	var err error

	// Get the list of filesystem objects in the helpers directory
	var filesystemObjects []os.FileInfo
	filesystemObjects, err = ioutil.ReadDir(*helpersDirPath)
	if err != nil {
		logger.Fatalln(err)
	}

	// Make sure that there are no directories in the helpers folder, and collect all of the file paths
	var filePaths []string
	for _, filesystemObject := range filesystemObjects {
		var filePath = filepath.Join(*helpersDirPath, filesystemObject.Name())
		if filesystemObject.IsDir() {
			logger.Fatalln(fmt.Sprintf("Sub-directories are not allowed in the \"templates\" folder.  Found sub-directory: %s", filePath))
		}

		filePaths = append(filePaths, filePath)
	}

	// Create template
	var tmpl = template.New("helpers")

	// Add options and functions
	tmpl = tmpl.Option("missingkey=error").Funcs(functionMap)

	// Parse helper files
	tmpl, err = tmpl.ParseFiles(filePaths...)
	if err != nil {
		logger.Fatalln(err)
	}

	return tmpl
}

func processTemplatesAndWriteToFilesystem(parentTemplate *template.Template, templatesDirPath *string, outputDirPath *string, parameters *map[interface{}]interface{}) {
	var err error

	// Get the list of filesystem objects in the templates directory
	var filesystemObjects []os.FileInfo
	filesystemObjects, err = ioutil.ReadDir(*templatesDirPath)
	if err != nil {
		logger.Fatalln(err)
	}

	// Make sure that there are no directories in the templates folder
	for _, filesystemObject := range filesystemObjects {
		if filesystemObject.IsDir() {
			logger.Fatalln(fmt.Sprintf("Sub-directories are not allowed in the \"templates\" folder.  Found sub-directory: %s", filepath.Join(*templatesDirPath, filesystemObject.Name())))
		}
	}

	// Delete and re-create the output directory in case it isn't empty or doesn't exist
	os.RemoveAll(*outputDirPath)
	os.MkdirAll(*outputDirPath, os.ModePerm)

	// Generate output from the templates
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()
		var filePath = filepath.Join(*templatesDirPath, fileName)
		var outputFilePath = filepath.Join(*outputDirPath, fileName)
		logger.Println(outputFilePath)

		// Generate the output
		var templateString = string(*readFileToBytes(&filePath))
		var generatedFileBytes = executeTemplate(parentTemplate, &fileName, &templateString, parameters)

		// Write the output to a file
		ioutil.WriteFile(outputFilePath, *generatedFileBytes, os.ModePerm)
	}
}

func executeTemplate(parentTemplate *template.Template, templateName *string, templateString *string, parameters *map[interface{}]interface{}) *[]byte {
	var err error

	// Create template object
	var tmpl *template.Template
	if parentTemplate != nil {
		tmpl = parentTemplate.New(*templateName)
	} else {
		tmpl = template.New(*templateName)
	}

	// Ensure failure on a failure to find input parameters
	tmpl = tmpl.Option("missingkey=error").Funcs(functionMap)

	// Parse the template
	tmpl, err = tmpl.Parse(*templateString)
	if err != nil {
		logger.Fatalln(err)
	}

	// Apply parameters to template
	var outputByteBuffer = new(bytes.Buffer)
	err = tmpl.Execute(outputByteBuffer, *parameters)
	if err != nil {
		logger.Fatalln(err)
	}

	// Convert bytes to a string
	var outputBytes = outputByteBuffer.Bytes()

	return &outputBytes
}

func readFileToBytes(filePath *string) *[]byte {
	var fileData, err = ioutil.ReadFile(*filePath)
	if err != nil {
		logger.Fatalln(err)
	}

	return &fileData
}

func yamlBytesToObject(yamlBytes *[]byte) *map[interface{}]interface{} {
	var result = make(map[interface{}]interface{})
	var err = yaml.UnmarshalStrict(*yamlBytes, &result)
	if err != nil {
		logger.Fatalln(err)
	}

	return &result
}

var functionMap = template.FuncMap{
	// Override the "index" function so it correctly fails the template generation on missing keys
	"index": func(values map[interface{}]interface{}, keys ...string) (interface{}, error) {
		if len(keys) == 0 {
			return values, nil
		}

		var currentMap = values
		var result interface{}
		for _, key := range keys {
			var ok bool
			result, ok = currentMap[key]
			if !ok {
				return nil, errors.New("Missing key: " + key)
			}

			// Try to assign the next map if the type is a map
			currentMap, ok = result.(map[interface{}]interface{})
			if !ok {
				// If the type is not a map, set this to nil so we don't reuse the old map
				currentMap = nil
			}
		}

		return result, nil
	},
}
