package subcommands

import (
	"bytes"
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

	// Define file locations
	//var dependenciesFilePath = filepath.Join(*packageDirPath, "dependencies.yaml")
	var interfaceFilePath = filepath.Join(*packageDirPath, "interface.yaml")
	//var packageFilePath = filepath.Join(*packageDirPath, "package.yaml")

	// Define directory locations
	//var dependenciesDirPath = filepath.Join(packageDirPath, "dependencies")
	var templatesDirPath = filepath.Join(*packageDirPath, "templates")

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
	var parameters = getParametersFromInterface(&interfaceFilePath, parametersFilePath)

	// Get the templates
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()
		var filePath = filepath.Join(templatesDirPath, fileName)
		var outputFilePath = filepath.Join(*outputDirPath, fileName)
		logger.Println(outputFilePath)

		// Generate the output
		var templateString = string(*readFileToBytes(&filePath))
		var generatedFileBytes = executeTemplate(&templateString, parameters)

		// Write the output to a file
		ioutil.WriteFile(outputFilePath, *generatedFileBytes, os.ModePerm)
	}

	logger.Println("Success!")

	return nil
}

func getParametersFromInterface(interfaceFilePath *string, userParametersFilePath *string) *map[interface{}]interface{} {
	// Get interface template as a string
	var templateString = string(*readFileToBytes(interfaceFilePath))

	// Get user parameters
	var userParameters = yamlBytesToObject(readFileToBytes(userParametersFilePath))

	// Generate interface from template and user parameters
	var interfaceBytes = executeTemplate(&templateString, userParameters)

	// Get values from generated interface
	var result = yamlBytesToObject(interfaceBytes)

	return result
}

func executeTemplate(templateString *string, parameters *map[interface{}]interface{}) *[]byte {
	var err error

	// Create template object
	var tmpl *template.Template
	tmpl, err = template.New("kpm-template").Option("missingkey=error").Parse(*templateString)
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
