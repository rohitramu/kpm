package templates

import (
	"fmt"
	"text/template"
)

type packageFuncFactory (func(tmpl *template.Template) any)

// GetPackageFuncMap returns the template functions which can be used only in the context of a particular template.
// If the template provided is nil, placeholder template functions are provided which return "Not implemented" errors.
func GetPackageFuncMap(tmpl *template.Template) map[string]any {
	return map[string]any{
		FuncNameInclude: getPackageFuncOrPlaceholder(tmpl, GetIncludeFunc),
	}
}

func getPackageFuncOrPlaceholder(tmpl *template.Template, fn packageFuncFactory) any {
	if tmpl == nil {
		return func(...any) (any, error) {
			return nil, fmt.Errorf("not implemented")
		}
	}

	return fn(tmpl)
}
