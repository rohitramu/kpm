package template_repository

type RepositoryInfo struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Location any    `yaml:"location"`
}
