package model

import (
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_repository"
)

type Config struct {
	LogLevel                    log.Level                            `yaml:"logLevel"`
	TemplatePackageRepositories []template_repository.RepositoryInfo `yaml:"repositories"`
}
