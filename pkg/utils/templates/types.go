package templates

import (
	"text/template"
)

// GenericMap is a generic map.
type GenericMap map[string]any

// TemplateSupplier is a function that supplies templates.
type TemplateSupplier func() *template.Template

// TemplateConsumer is a function that consumes templates.
type TemplateConsumer func(tmpl *template.Template)

// PackageInfo is the structure of the package definition's yaml file.
type PackageInfo struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

// PackageDefinition contains the information required to execute a template package.
type PackageDefinition struct {
	PackageInfo *PackageInfo `yaml:"package" json:"package"`
	Parameters  *GenericMap  `yaml:"parameters" json:"parameters"`
}
