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
	var err error

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

	// Define directory locations
	//var dependenciesDirPath = filepath.Join(packageDirPath, dependenciesDirName)
	var templatesDirPath = filepath.Join(*packageDirPath, templatesDirName)

	// Get the list of filesystem objects in the templates directory
	var filesystemObjects []os.FileInfo
	filesystemObjects, err = ioutil.ReadDir(templatesDirPath)
	if err != nil {
		logger.Fatalln(err)
	}

	// Make sure that there are no directories in the templates folder
	for _, filesystemObject := range filesystemObjects {
		if filesystemObject.IsDir() {
			logger.Fatalln(fmt.Sprintf("Sub-directories are not allowed in the \"templates\" folder.  Found sub-directory: %s", filepath.Join(templatesDirPath, filesystemObject.Name())))
		}
	}

	// Delete the output directory in case it isn't empty
	os.RemoveAll(*outputDirPath)

	// Create the output directory if it didn't already exist
	os.MkdirAll(*outputDirPath, os.ModePerm)

	// Read the interface file and apply the template
	var interfaceTemplateString = string(*readFileToBytes(&interfaceFilePath))
	var interfaceBytes = executeTemplate(&interfaceFileName, &interfaceTemplateString, yamlBytesToObject(readFileToBytes(parametersFilePath)))
	var parameters = yamlBytesToObject(interfaceBytes)

	// Get the templates
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()
		var filePath = filepath.Join(templatesDirPath, fileName)
		var outputFilePath = filepath.Join(*outputDirPath, fileName)
		logger.Println(outputFilePath)

		// Generate the output
		var templateString = string(*readFileToBytes(&filePath))
		var generatedFileBytes = executeTemplate(&fileName, &templateString, parameters)

		// Write the output to a file
		ioutil.WriteFile(outputFilePath, *generatedFileBytes, os.ModePerm)
	}

	logger.Println("Success!")

	return nil
}

func executeTemplate(templateName *string, templateString *string, parameters *map[interface{}]interface{}) *[]byte {
	var err error

	// Create template object
	var tmpl *template.Template
	tmpl, err = template.New(*templateName).Option("missingkey=error").Funcs(functionMap).Parse(*templateString)
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
