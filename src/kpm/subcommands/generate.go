package subcommands

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"

	"./utils"
)

// Logger
var logger = utils.NewLogger()

// Root template
var rootTemplate = getRootTemplate()

// A generic map
type genericMap map[interface{}]interface{}

// Define top-level objects in template inputs
const (
	templateFieldPackage = "package"
	templateFieldValues  = "values"
)

// The structure of the package definition's yaml file
type packageInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// The structure of a dependency definition's yaml file
type dependencyPackageInfo struct {
	PackageInfo packageInfo `yaml:"package"`
	Values      genericMap  `yaml:"values"`
}

// GenerateCmd creates Kubernetes configuration files from the
// given template package directory and parameters file, and then
// writes them to the given output directory
func GenerateCmd(packageDirPathArg *string, parametersFilePathArg *string, outputDirPathArg *string) error {
	// Define well-known file names
	const (
		interfaceFileName  = "interface.yaml"
		packageFileName    = "package.yaml"
		parametersFileName = "parameters.yaml"
	)

	// Define well-known directory names
	const (
		dependenciesDirName = "dependencies"
		templatesDirName    = "templates"
		helpersDirName      = "helpers"
		outputDirName       = "_generated"
	)

	// Resolve package path
	var packageDirPath = utils.GetAbsolutePathOrDefault(packageDirPathArg, utils.GetCurrentWorkingDir())

	// Resolve parameters file path
	var parametersFilePath = utils.GetAbsolutePathOrDefault(parametersFilePathArg, filepath.Join(packageDirPath, parametersFileName))

	// Resolve output directory
	var outputDirPath = utils.GetAbsolutePathOrDefault(outputDirPathArg, filepath.Join(packageDirPath, outputDirName))

	// Log resolved arguments
	logger.Verbose.Println(fmt.Sprintf("Package directory: %s", packageDirPath))
	logger.Verbose.Println(fmt.Sprintf("Parameters file:   %s", parametersFilePath))
	logger.Verbose.Println(fmt.Sprintf("Output directory:  %s", outputDirPath))

	// Define file locations
	var interfaceFilePath = filepath.Join(packageDirPath, interfaceFileName)
	var packageFilePath = filepath.Join(packageDirPath, packageFileName)

	// Define directory locations
	var dependenciesDirPath = filepath.Join(packageDirPath, dependenciesDirName)
	var templatesDirPath = filepath.Join(packageDirPath, templatesDirName)
	var helpersDirPath = filepath.Join(packageDirPath, helpersDirName)

	// Get template input values by applying parameters to interface
	var templateInput = getTemplateInput(packageFilePath, interfaceFilePath, parametersFilePath)

	// Get helpers template
	var helpersTemplate = getHelpersTemplate(rootTemplate, helpersDirPath)

	// Generate output files from dependencies
	processDependenciesAndWriteToFilesystem(dependenciesDirPath, outputDirPath, helpersTemplate, templateInput)

	// Generate output files and write them to the output directory
	processTemplatesAndWriteToFilesystem(templatesDirPath, outputDirPath, helpersTemplate, templateInput)

	return nil
}

// +----------------------+
// | Process dependencies |
// +----------------------+

func processDependenciesAndWriteToFilesystem(dependenciesDirPath string, outputDirPath string, parentTemplate *template.Template, templateInput *genericMap) {

}

// +-------------------+
// | Process templates |
// +-------------------+

func processTemplatesAndWriteToFilesystem(templatesDirPath string, outputDirPath string, parentTemplate *template.Template, templateInput *genericMap) {
	var err error

	// Get the list of filesystem objects in the templates directory
	var filesystemObjects []os.FileInfo
	filesystemObjects, err = ioutil.ReadDir(templatesDirPath)
	if err != nil {
		logger.Error.Panicln(err)
	}

	// Make sure that there are no directories in the templates folder before starting the generation
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()
		if filesystemObject.IsDir() {
			logger.Error.Fatalln(fmt.Sprintf("Directories are not allowed inside the \"%s\" folder.  Found sub-directory: %s", templatesDirPath, fileName))
		}
	}

	// Delete and re-create the output directory in case it isn't empty or doesn't exist
	os.RemoveAll(outputDirPath)
	os.MkdirAll(outputDirPath, os.ModePerm)

	// Generate output from the templates
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()
		var filePath = filepath.Join(templatesDirPath, fileName)
		var outputFilePath = filepath.Join(outputDirPath, fileName)
		logger.Verbose.Println(fmt.Sprintf("Generating: %s", outputFilePath))

		// Generate the output
		var templateString = readFileToString(filePath)
		var generatedFileBytes = executeTemplate(parentTemplate, fileName, templateString, templateInput)

		// Write the output to a file
		ioutil.WriteFile(outputFilePath, generatedFileBytes, os.ModeAppend|os.ModePerm)
	}

	// Print status
	logger.Info.Println(fmt.Sprintf("SUCCESS: %s", outputDirPath))
}

func executeTemplate(parentTemplate *template.Template, templateName string, templateString string, data interface{}) []byte {
	var err error

	// Create template object
	var tmpl *template.Template
	if parentTemplate == nil {
		logger.Error.Panicln("A parent template must be provided")
	}

	// Create a new template which references the parent template
	tmpl = parentTemplate.New(templateName)

	// Parse the template
	tmpl, err = tmpl.Parse(templateString)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	// Apply data to template
	var outputByteBuffer = new(bytes.Buffer)
	err = tmpl.Execute(outputByteBuffer, data)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	// Convert bytes to a string
	var outputBytes = outputByteBuffer.Bytes()

	return outputBytes
}

var functionMap = template.FuncMap{
	// Override the "index" function so it correctly fails the template generation on missing keys
	"index": func(data genericMap, keys ...interface{}) (interface{}, error) {
		if len(keys) == 0 {
			return data, nil
		}

		var currentMap = data
		var result interface{}
		for _, key := range keys {
			var ok bool
			result, ok = currentMap[key]
			if !ok {
				var keyName string
				keyName, ok = key.(string)
				var message string
				if !ok {
					message = fmt.Sprintf("Missing key of type: %s", reflect.TypeOf(key))
				} else {
					message = fmt.Sprintf("Missing key: %s", keyName)
				}
				return nil, errors.New(message)
			}

			// Try to assign the next map if the type is a map
			currentMap, ok = result.(genericMap)
			if !ok {
				// If the type is not a map, set this to nil so we don't reuse the old map
				currentMap = nil
			}
		}

		return result, nil
	},
}

// +------------+
// | Get values |
// +------------+

func getTemplateInput(packageFilePath string, interfaceFilePath string, parametersFilePath string) *genericMap {
	// Get package info
	var packageInfo = getPackageInfo(packageFilePath)

	// Generate the values by populating the interface template with parameters
	var values = getValuesFromInterface(interfaceFilePath, parametersFilePath)

	// Add top-level objects
	var result = genericMap{}
	result[templateFieldPackage] = packageInfo
	result[templateFieldValues] = values

	return &result
}

func getPackageInfo(packageInfoFilePath string) *packageInfo {
	// Get file content
	var yamlBytes = readFileToBytes(packageInfoFilePath)

	// Get PackageInfo object from file content
	var result = yamlBytesToPackageInfo(yamlBytes)

	return result
}

func getValuesFromInterface(interfaceFilePath string, parametersFilePath string) *genericMap {
	// Get interface template as a string
	var templateString = readFileToString(interfaceFilePath)

	// Get parameters
	var parameters = yamlBytesToMap(readFileToBytes(parametersFilePath))

	// Generate values yaml file (in-memory) by applying parameters to interface
	var interfaceFileName = filepath.Base(interfaceFilePath)
	var interfaceBytes = executeTemplate(rootTemplate, interfaceFileName, templateString, parameters)

	// Get values object from generated values yaml file
	var result = yamlBytesToMap(interfaceBytes)

	return result
}

func getHelpersTemplate(parentTemplate *template.Template, helpersDirPath string) *template.Template {
	var err error

	// Get the list of filesystem objects in the helpers directory
	var filesystemObjects []os.FileInfo
	filesystemObjects, err = ioutil.ReadDir(helpersDirPath)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	// Make sure that there are no directories in the helpers folder, and collect all of the file paths
	var filePaths []string
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()
		var filePath = filepath.Join(helpersDirPath, fileName)
		if filesystemObject.IsDir() {
			logger.Error.Fatalln(fmt.Sprintf("Sub-directories are not allowed in the \"%s\" folder.  Found sub-directory: %s", helpersDirPath, fileName))
		}

		filePaths = append(filePaths, filePath)
	}

	// Create template
	var tmpl = template.New("helpers")

	// Add options and functions
	tmpl = tmpl.Option("missingkey=error").Funcs(functionMap).Funcs(sprig.TxtFuncMap())

	// Parse helper files
	tmpl, err = tmpl.ParseFiles(filePaths...)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	return tmpl
}

func getRootTemplate() *template.Template {
	// Create template
	var tmpl = template.New("root")

	// Make sure template execution fails if a key is missing
	tmpl = tmpl.Option("missingkey=error")

	// Add sprig functions
	tmpl.Funcs(sprig.TxtFuncMap())

	// Add custom functions
	tmpl.Funcs(functionMap)

	return tmpl
}

// +-----------+
// | Read file |
// +-----------+

func readFileToString(filePath string) string {
	var result = string(readFileToBytes(filePath))

	return result
}

func readFileToBytes(filePath string) []byte {
	var fileData, err = ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	return fileData
}

// +------------------------------+
// | Convert yaml bytes to object |
// +------------------------------+

func yamlBytesToMap(yamlBytes []byte) *genericMap {
	// NOTE: ALWAYS pass "UnmarshalStrict()" a pointer rather than a real value
	var result = &genericMap{}
	var err = yaml.UnmarshalStrict(yamlBytes, result)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	return result
}

func yamlBytesToPackageInfo(packageInfoBytes []byte) *packageInfo {
	// NOTE: ALWAYS pass "UnmarshalStrict()" a pointer rather than a real value
	var result = &packageInfo{}
	var err = yaml.UnmarshalStrict(packageInfoBytes, result)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	return result
}
