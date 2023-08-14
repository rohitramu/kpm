package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"

	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/types"
)

// NewRootTemplate returns a new root template with options and functions provided.
func NewRootTemplate() *template.Template {
	// Create template
	var tmpl = template.New("root")

	// Make sure template execution fails if a key is missing
	tmpl = tmpl.Option("missingkey=error")

	// Add sprig functions
	tmpl = tmpl.Funcs(sprig.TxtFuncMap())

	// Add global functions
	tmpl = tmpl.Funcs(GetGlobalFuncMap())

	// Add placeholders for package-specific functions
	tmpl = tmpl.Funcs(GetPackageFuncMap(nil))

	return tmpl
}

// AddPackageSpecificTemplateFunctions adds the package-specific template functions for the given template.
func AddPackageSpecificTemplateFunctions(tmpl *template.Template) *template.Template {
	if tmpl == nil {
		log.Panicf("Template cannot be nil")
	}

	return tmpl.Funcs(GetPackageFuncMap(tmpl))
}

// GetTemplateFromFile returns a new template object given a template file.
func GetTemplateFromFile(parentTemplate *template.Template, templateName string, filePath string) (*template.Template, error) {
	var err error

	// Create template
	var tmpl *template.Template
	if parentTemplate != nil {
		tmpl = parentTemplate.New(templateName)
	} else {
		tmpl = template.New(templateName)
	}

	// Get template file as string
	var templateString string
	templateString, err = files.ReadString(filePath)
	if err != nil {
		return nil, err
	}

	// Parse template
	tmpl, err = tmpl.Parse(templateString)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// GetTemplatesFromDir returns an array containing all of the templates found in the given directory.
func GetTemplatesFromDir(parentTemplate *template.Template, templatesDirPath string) ([]*template.Template, error) {
	var err error

	var templates []*template.Template
	err = visitTemplatesFromDir(templatesDirPath, func() *template.Template {
		// Use the same parent template each time
		return parentTemplate
	}, func(tmpl *template.Template) {
		// Add the template to the array
		templates = append(templates, tmpl)
	})

	if err != nil {
		return nil, err
	}

	return templates, nil
}

// ChainTemplatesFromDir returns a single template which contains all of the templates that were found in the given directory.
func ChainTemplatesFromDir(parentTemplate *template.Template, templatesDirPath string) (*template.Template, int, error) {
	var err error

	var currentTemplate = parentTemplate
	var numTemplates = 0
	err = visitTemplatesFromDir(templatesDirPath, func() *template.Template {
		// Use the current template as the parent
		return currentTemplate
	}, func(nextTemplate *template.Template) {
		// Increment template count
		numTemplates++

		// Set the next template as current
		currentTemplate = nextTemplate
	})
	if err != nil {
		return nil, 0, err
	}

	return currentTemplate, numTemplates, nil
}

// ExecuteNamedTemplate executes a named template that can be found in the provided template.
func ExecuteNamedTemplate(tmpl *template.Template, templateName string, values interface{}) ([]byte, error) {
	var err error

	// Check that the parent template is not nil
	if tmpl == nil {
		log.Panicf("The provided template cannot be nil")
	}

	// Check that the template name is not empty
	if templateName == "" {
		return nil, fmt.Errorf("Template name cannot be empty: %s", tmpl.Name())
	}

	// Get the named template
	var namedTemplate = tmpl.Lookup(templateName)
	if namedTemplate == nil {
		return nil, fmt.Errorf("Failed to find named template: %s", templateName)
	}

	// Execute the named template with the provided values
	var result []byte
	result, err = ExecuteTemplate(namedTemplate, values)

	return result, err
}

// ExecuteTemplate executes a template given the template object and the values.
func ExecuteTemplate(tmpl *template.Template, values interface{}) ([]byte, error) {
	var err error

	// Create template object
	if tmpl == nil {
		return nil, fmt.Errorf("The template to execute cannot be nil")
	}

	if values == nil {
		return nil, fmt.Errorf("The values to execute the template with cannot be nil")
	}

	// Apply values to template
	log.Debugf("Executing template: %s", tmpl.Name())
	var outputByteBuffer = new(bytes.Buffer)
	err = tmpl.Execute(outputByteBuffer, values)
	if err != nil {
		return nil, err
	}

	// Convert bytes to a string
	var outputBytes = outputByteBuffer.Bytes()

	return outputBytes, nil
}

// visitTemplatesFromDir visits each template found in the given directory, sets the parent using the given "getParentTemplate" function
// and then consumes the template using the given "consumeTemplate" function.
func visitTemplatesFromDir(templatesDirPath string, getParentTemplate types.TemplateSupplier, consumeTemplate types.TemplateConsumer) error {
	var err error

	// Get the list of filesystem objects in the helpers directory
	var filesystemObjects []os.FileInfo
	filesystemObjects, err = ioutil.ReadDir(templatesDirPath)
	if err != nil {
		return err
	}

	// Parse all templates in the given directory, ignoring sub-directories
	log.Debugf("Parsing templates in directory: %s", templatesDirPath)
	for _, filesystemObject := range filesystemObjects {
		var fileName = filesystemObject.Name()

		// Ignore directories
		if filesystemObject.IsDir() {
			log.Warningf("Ignoring sub-directory: %s", fileName)
			continue
		}

		log.Debugf("Parsing template: %s", fileName)

		// Create a template object from the file
		var filePath = filepath.Join(templatesDirPath, fileName)
		var tmpl *template.Template
		tmpl, err = GetTemplateFromFile(getParentTemplate(), fileName, filePath)
		if err != nil {
			return err
		}

		// Consume template
		log.Debugf("Consuming template: %s", tmpl.Name())
		consumeTemplate(tmpl)
	}

	return nil
}
