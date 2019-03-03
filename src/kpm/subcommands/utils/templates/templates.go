package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"

	"../files"
	"../logger"
	"../templatefuncs"
	"../types"
)

// GetRootTemplate returns a new root template with options and functions provided.
func GetRootTemplate() *template.Template {
	// Create template
	var tmpl = template.New("root")

	// Make sure template execution fails if a key is missing
	tmpl = tmpl.Option("missingkey=error")

	// Add sprig functions
	tmpl.Funcs(sprig.TxtFuncMap())

	// Add custom functions
	tmpl.Funcs(template.FuncMap{
		// Override the "index" function so it correctly fails the template generation on missing keys
		"index": templatefuncs.Index,
	})

	return tmpl
}

// GetTemplateFromFile returns a new template object given a template file.
func GetTemplateFromFile(parentTemplate *template.Template, templateName string, filePath string) *template.Template {
	var err error

	// Create template
	var tmpl *template.Template
	if parentTemplate != nil {
		tmpl = parentTemplate.New(templateName)
	} else {
		tmpl = template.New(templateName)
	}

	// Get template file as string
	var templateString = files.ReadFileToString(filePath)

	// Parse template
	tmpl, err = tmpl.Parse(templateString)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	return tmpl
}

// GetTemplatesFromDir returns an array containing all of the templates found in the given directory.
func GetTemplatesFromDir(parentTemplate *template.Template, templatesDirPath string) []*template.Template {
	var templates []*template.Template
	visitTemplatesFromDir(templatesDirPath, func() *template.Template {
		// Use the same parent template each time
		return parentTemplate
	}, func(tmpl *template.Template) {
		// Add the template to the array
		templates = append(templates, tmpl)
	})

	return templates
}

// ChainTemplatesFromDir returns a single template which contains all of the templates that were found in the given directory.
func ChainTemplatesFromDir(parentTemplate *template.Template, templatesDirPath string) (*template.Template, int) {
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

// ExecuteTemplate executes a template given the template object and the values.
func ExecuteTemplate(tmpl *template.Template, values *types.GenericMap) []byte {
	var err error

	// Create template object
	if tmpl == nil {
		logger.Default.Error.Panicln("The template to execute cannot be nil")
	}

	// Apply values to template
	var outputByteBuffer = new(bytes.Buffer)
	err = tmpl.Execute(outputByteBuffer, (*types.GenericMap)(values))
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Convert bytes to a string
	var outputBytes = outputByteBuffer.Bytes()

	return outputBytes
}

// visitTemplatesFromDir visits each template found in the given directory, sets the parent using the given "getParentTemplate" function
// and then consumes the template using the given "consumeTemplate" function.  Finally, this function returns a count of the number of
// templates that were visited.
func visitTemplatesFromDir(templatesDirPath string, getParentTemplate types.TemplateSupplier, consumeTemplate types.TemplateConsumer) int {
	var err error

	// Get the list of filesystem objects in the helpers directory
	var filesystemObjects []os.FileInfo
	filesystemObjects, err = ioutil.ReadDir(templatesDirPath)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}

	// Parse all templates in the given directory, ignoring sub-directories
	logger.Default.Info.Println(fmt.Sprintf("Parsing templates in directory: %s", templatesDirPath))
	var numTemplates = 0
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()

		// Ignore directories
		if filesystemObject.IsDir() {
			logger.Default.Warning.Println(fmt.Sprintf("Ignoring sub-directory: %s", fileName))
		} else {
			logger.Default.Verbose.Println(fmt.Sprintf("Parsing template: %s", fileName))

			// Create a template object from the file
			var filePath = filepath.Join(templatesDirPath, fileName)
			var tmpl = GetTemplateFromFile(getParentTemplate(), fileName, filePath)

			// Consume template
			logger.Default.Verbose.Println(fmt.Sprintf("Consuming template: %s", tmpl.Name()))
			consumeTemplate(tmpl)

			// Increment template count
			numTemplates++
		}
	}

	return numTemplates
}
