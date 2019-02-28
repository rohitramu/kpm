package utils

import (
	"text/template"
)

// GenericMap is a generic map
type GenericMap map[interface{}]interface{}

// TemplateSupplier is a function that supplies templates
type TemplateSupplier func() *template.Template

// TemplateConsumer is a function that consumes templates
type TemplateConsumer func(tmpl *template.Template)
