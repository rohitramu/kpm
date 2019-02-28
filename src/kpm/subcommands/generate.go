package subcommands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"

	"./templatefuncs"
	"./utils"
)

// Logger
var logger = utils.NewLogger()

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
	PackageInfo packageInfo      `yaml:"package"`
	Values      utils.GenericMap `yaml:"values"`
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
		outputDirName       = ".kpm_generated"
	)

	// Resolve paths
	var (
		workingDir          = utils.GetCurrentWorkingDir()
		packageDirPath      = utils.GetAbsolutePathOrDefault(packageDirPathArg, workingDir)
		packageName         = filepath.Base(packageDirPath)
		outputDirParentPath = utils.GetAbsolutePathOrDefault(outputDirPathArg, filepath.Join(workingDir, outputDirName))
		outputDirPath       = filepath.Join(outputDirParentPath, packageName)
		parametersFilePath  = utils.GetAbsolutePathOrDefault(parametersFilePathArg, filepath.Join(packageDirPath, parametersFileName))
	)

	// Log resolved paths
	logger.Verbose.Println("====")
	logger.Verbose.Println(fmt.Sprintf("Package directory: %s", packageDirPath))
	logger.Verbose.Println(fmt.Sprintf("Parameters file:   %s", parametersFilePath))
	logger.Verbose.Println(fmt.Sprintf("Output directory:  %s", outputDirPath))
	logger.Verbose.Println("====")

	// Define file locations
	var (
		interfaceFilePath = filepath.Join(packageDirPath, interfaceFileName)
		packageFilePath   = filepath.Join(packageDirPath, packageFileName)
	)

	// Define directory locations
	var (
		dependenciesDirPath = filepath.Join(packageDirPath, dependenciesDirName)
		templatesDirPath    = filepath.Join(packageDirPath, templatesDirName)
		helpersDirPath      = filepath.Join(packageDirPath, helpersDirName)
	)

	// Get template from helpers
	var helpersTemplate, numHelpers = chainTemplatesFromDir(getRootTemplate(), helpersDirPath)
	logger.Verbose.Println(fmt.Sprintf("Found %d helper template(s) in directory: %s", numHelpers, helpersDirPath))

	// Get template input values by applying parameters to interface
	var templateInput = getTemplateInput(helpersTemplate, packageFilePath, interfaceFilePath, parametersFilePath)

	// Generate output files from dependencies
	processDependenciesAndWriteToFilesystem(dependenciesDirPath, outputDirPath, helpersTemplate, templateInput)

	// Generate output files and write them to the output directory
	var numProcessedTemplates = processTemplatesAndWriteToFilesystem(helpersTemplate, templatesDirPath, templateInput, outputDirPath)
	logger.Verbose.Println(fmt.Sprintf("Processed %d template(s) in directory: %s", numProcessedTemplates, templatesDirPath))

	// Print status
	logger.Info.Println(fmt.Sprintf("SUCCESS - Generated output in directory: %s", outputDirPath))

	return nil
}

// +----------------------+
// | Process dependencies |
// +----------------------+

func processDependenciesAndWriteToFilesystem(dependenciesDirPath string, outputDirPath string, parentTemplate *template.Template, templateInput *utils.GenericMap) {

}

// +-------------------+
// | Process templates |
// +-------------------+

func processTemplatesAndWriteToFilesystem(parentTemplate *template.Template, templatesDirPath string, templateInput *utils.GenericMap, outputDirPath string) int {
	// Delete and re-create the output directory in case it isn't empty or doesn't exist
	os.RemoveAll(outputDirPath)
	os.MkdirAll(outputDirPath, os.ModePerm)

	var numTemplates = visitTemplatesFromDir(templatesDirPath, func() *template.Template {
		// Use the given parent template
		return parentTemplate
	}, func(tmpl *template.Template) {
		// Generate output from each template
		var generatedFileBytes = executeTemplate(tmpl, templateInput)

		// Write the output to a file
		var outputFilePath = filepath.Join(outputDirPath, tmpl.Name())
		ioutil.WriteFile(outputFilePath, generatedFileBytes, os.ModeAppend|os.ModePerm)
	})

	return numTemplates
}

// +------------+
// | Get values |
// +------------+

func getTemplateInput(parentTemplate *template.Template, packageFilePath string, interfaceFilePath string, parametersFilePath string) *utils.GenericMap {
	// Get package info
	var packageInfo = getPackageInfo(packageFilePath)

	// Generate the values by populating the interface template with parameters
	var values = getValuesFromInterface(parentTemplate, interfaceFilePath, parametersFilePath)

	// Add top-level objects
	var result = utils.GenericMap{}
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

func getValuesFromInterface(parentTemplate *template.Template, interfaceFilePath string, parametersFilePath string) *utils.GenericMap {
	var err error

	// Create template object from interface file
	var templateName = filepath.Base(interfaceFilePath)
	var tmpl = getTemplateFromFile(parentTemplate, templateName, interfaceFilePath)

	// Get parameters file content as bytes
	var parametersFileBytes []byte
	parametersFileBytes, err = ioutil.ReadFile(parametersFilePath)
	if err != nil {
		logger.Warning.Println(fmt.Sprintf("Failed to read parameters file: %s", err))
		parametersFileBytes = []byte{}
	}

	// Get parameters
	var parameters = yamlBytesToMap(parametersFileBytes)

	// Generate values by applying parameters to interface
	var interfaceBytes = executeTemplate(tmpl, parameters)

	// Get values object from generated values yaml file
	var result = yamlBytesToMap(interfaceBytes)

	return result
}

// +-----------+
// | Templates |
// +-----------+

func getRootTemplate() *template.Template {
	// Create template
	var tmpl = template.New("root")

	// Make sure template execution fails if a key is missing
	tmpl = tmpl.Option("missingkey=error")

	// Add sprig functions
	tmpl.Funcs(sprig.TxtFuncMap())

	// Add custom functions
	tmpl.Funcs(templateFunctions)

	return tmpl
}

func getTemplateFromFile(parentTemplate *template.Template, templateName string, filePath string) *template.Template {
	var err error

	// Create template
	var tmpl *template.Template
	if parentTemplate != nil {
		tmpl = parentTemplate.New(templateName)
	} else {
		tmpl = template.New(templateName)
	}

	// Get template file as string
	var templateString = readFileToString(filePath)

	// Parse template
	tmpl, err = tmpl.Parse(templateString)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	return tmpl
}

func getTemplatesFromDir(parentTemplate *template.Template, templatesDirPath string) ([]*template.Template, int) {
	var templates []*template.Template
	var numTemplates = visitTemplatesFromDir(templatesDirPath, func() *template.Template {
		// Use the same parent template each time
		return parentTemplate
	}, func(tmpl *template.Template) {
		// Add the template to the array
		templates = append(templates, tmpl)
	})

	return templates, numTemplates
}

func chainTemplatesFromDir(parentTemplate *template.Template, templatesDirPath string) (*template.Template, int) {
	var currentTemplate = parentTemplate
	var numTemplates = visitTemplatesFromDir(templatesDirPath, func() *template.Template {
		// Use the current template as the parent
		return currentTemplate
	}, func(nextTemplate *template.Template) {
		// Set the next template as current
		currentTemplate = nextTemplate
	})

	return currentTemplate, numTemplates
}

func visitTemplatesFromDir(templatesDirPath string, getParentTemplate utils.TemplateSupplier, consumeTemplate utils.TemplateConsumer) int {
	var err error

	// Get the list of filesystem objects in the helpers directory
	var filesystemObjects []os.FileInfo
	filesystemObjects, err = ioutil.ReadDir(templatesDirPath)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	// Parse all templates in the given directory, ignoring sub-directories
	logger.Info.Println(fmt.Sprintf("Parsing templates in directory: %s", templatesDirPath))
	var numTemplates = 0
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()

		// Ignore directories
		if filesystemObject.IsDir() {
			logger.Warning.Println(fmt.Sprintf("Ignoring sub-directory: %s", fileName))
		} else {
			logger.Verbose.Println(fmt.Sprintf("Parsing template: %s", fileName))

			// Create a template object from the file
			var filePath = filepath.Join(templatesDirPath, fileName)
			var tmpl = getTemplateFromFile(getParentTemplate(), fileName, filePath)

			// Consume template
			logger.Verbose.Println(fmt.Sprintf("Consuming template: %s", tmpl.Name()))
			consumeTemplate(tmpl)

			// Increment template count
			numTemplates++
		}
	}

	return numTemplates
}

func executeTemplate(tmpl *template.Template, data interface{}) []byte {
	var err error

	// Create template object
	if tmpl == nil {
		logger.Error.Panicln("The template to execute cannot be nil")
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

var templateFunctions = template.FuncMap{
	// Override the "index" function so it correctly fails the template generation on missing keys
	"index": templatefuncs.Index,
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

func yamlBytesToMap(yamlBytes []byte) *utils.GenericMap {
	// NOTE: ALWAYS pass "UnmarshalStrict()" a pointer rather than a real value
	var result = &utils.GenericMap{}
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
