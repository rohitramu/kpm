package templates

import (
	"fmt"
	"text/template"

	"github.com/rohitramu/kpm/subcommands/utils/log"
)

// FuncNameInclude is the name of the "include" template function.
const FuncNameInclude = "include"

// GetIncludeFunc creates a new instance of the Include function, which allows helper templates to be executed so their output can be used in other functions.
func GetIncludeFunc(tmpl *template.Template) PackageFunc {
	if tmpl == nil {
		log.Panic("Template cannot be nil")
	}

	return func(parameters ...interface{}) (interface{}, error) {
		var err error
		var ok bool

		// Make sure we have parameters
		if parameters == nil {
			return nil, fmt.Errorf("Parameters must passed to \"%s\"", FuncNameInclude)
		}

		// Make sure we have at least the name of the template to include
		if len(parameters) < 2 {
			return nil, fmt.Errorf("At 1 parameter (the name of the template) must be provided to \"%s\"", FuncNameInclude)
		}

		// Make sure we don't have any extra parameters
		if len(parameters) > 2 {
			return nil, fmt.Errorf("Only 2 parameters can be provided to \"%s\"", FuncNameInclude)
		}

		// Get the template name
		var templateName string
		templateName, ok = parameters[0].(string)
		if !ok {
			return nil, fmt.Errorf("The first parameter to the \"%s\" function must be a string (the name of the template to include)", FuncNameInclude)
		}

		// Get the value to pass to the template if we were given any
		var values interface{}
		if len(parameters) > 1 {
			values = parameters[1]
		}

		// Execute the named template
		var resultBytes []byte
		resultBytes, err = ExecuteNamedTemplate(tmpl, templateName, values)
		if err != nil {
			return nil, fmt.Errorf("Failed to execute named template \"%s\": %s", templateName, err)
		}

		// Get the result as a string
		var result = string(resultBytes)

		return result, nil
	}
}
