package templates

import (
	"fmt"
	"text/template"
)

// PackageFunc represents a template function which is package-specific.
type PackageFunc func(...interface{}) (interface{}, error)

type packageFuncFactory (func(tmpl *template.Template) PackageFunc)

// GetPackageFuncMap returns the template functions which can be used only in the context of a particular template.
// If the template provided is nil, placeholder template functions are provided which return "Not implemented" errors.
func GetPackageFuncMap(tmpl *template.Template) map[string]interface{} {
	return map[string]interface{}{
		FuncNameInclude: getPackageFuncOrPlaceholder(tmpl, GetIncludeFunc),
	}
}

func getPackageFuncOrPlaceholder(tmpl *template.Template, fn packageFuncFactory) PackageFunc {
	if tmpl == nil {
		return func(...interface{}) (interface{}, error) {
			return nil, fmt.Errorf("Not implemented")
		}
	}

	return fn(tmpl)
}
