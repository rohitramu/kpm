package types

import (
	"text/template"
)

// GenericMap is a generic map
type GenericMap map[string]interface{}

// TemplateSupplier is a function that supplies templates
type TemplateSupplier func() *template.Template

// TemplateConsumer is a function that consumes templates
type TemplateConsumer func(tmpl *template.Template)

// PackageInfo is the structure of the package definition's yaml file
type PackageInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// DependencyPackageInfo is the structure of a dependency definition's yaml file
type DependencyPackageInfo struct {
	PackageInfo PackageInfo `yaml:"package"`
	Values      GenericMap  `yaml:"values"`
}
