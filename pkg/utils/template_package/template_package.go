package template_package

// PackageInfo is the structure of the package definition's yaml file.
type PackageInfo struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

// PackageDefinition contains the information required to execute a template package.
type PackageDefinition struct {
	PackageInfo *PackageInfo    `yaml:"package" json:"package"`
	Parameters  *map[string]any `yaml:"parameters" json:"parameters"`
}

func (p PackageInfo) String() string {
	return GetPackageFullName(p.Name, p.Version)
}
