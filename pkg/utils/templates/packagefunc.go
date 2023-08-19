package templates

import (
	"fmt"
	"text/template"
)

type packageFuncFactory (func(tmpl *template.Template) interface{})

// GetPackageFuncMap returns the template functions which can be used only in the context of a particular template.
// If the template provided is nil, placeholder template functions are provided which return "Not implemented" errors.
func GetPackageFuncMap(tmpl *template.Template) map[string]interface{} {
	return map[string]interface{}{
		FuncNameInclude: getPackageFuncOrPlaceholder(tmpl, GetIncludeFunc),
	}
}

func getPackageFuncOrPlaceholder(tmpl *template.Template, fn packageFuncFactory) interface{} {
	if tmpl == nil {
		return func(...interface{}) (interface{}, error) {
			return nil, fmt.Errorf("not implemented")
		}
	}

	return fn(tmpl)
}
