package templates

import (
	"fmt"
	"text/template"

	"github.com/rohitramu/kpm/subcommands/utils/log"
)

// FuncNameInclude is the name of the "include" template function.
const FuncNameInclude = "include"

// GetIncludeFunc creates a new instance of the Include function, which allows helper templates to be executed so their output can be used in other functions.
func GetIncludeFunc(tmpl *template.Template) interface{} {
	if tmpl == nil {
		log.Panic("Template cannot be nil")
	}

	return func(templateName string, data interface{}) (string, error) {
		var err error

		// Execute the named template
		var resultBytes []byte
		resultBytes, err = ExecuteNamedTemplate(tmpl, templateName, data)
		if err != nil {
			return "", fmt.Errorf("Failed to execute named template \"%s\":\n%s", templateName, err)
		}

		// Get the result as a string
		var result = string(resultBytes)

		return result, nil
	}
}
