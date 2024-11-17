package templates

import (
	"text/template"
)

// TemplateSupplier is a function that supplies templates.
type TemplateSupplier func() *template.Template

// TemplateConsumer is a function that consumes templates.
type TemplateConsumer func(tmpl *template.Template)
