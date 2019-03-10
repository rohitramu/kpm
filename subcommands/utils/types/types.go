package types

import (
	"text/template"
)

// GenericMap is a generic map.
type GenericMap map[string]interface{}

// TemplateSupplier is a function that supplies templates.
type TemplateSupplier func() *template.Template

// TemplateConsumer is a function that consumes templates.
type TemplateConsumer func(tmpl *template.Template)

// PackageInfo is the structure of the package definition's yaml file.
type PackageInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// PackageDefinition contains the information required to execute a template package.
type PackageDefinition struct {
	Package    *PackageInfo `yaml:"package"`
	Parameters *GenericMap  `yaml:"parameters"`
}
