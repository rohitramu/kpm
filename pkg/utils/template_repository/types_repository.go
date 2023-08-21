package template_repository

import (
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

type Repository interface {
	GetName() string
	GetType() string
	Packages() ([]templates.PackageInfo, error)
	PackageVersions() ([]string, error)
	Push(templates.PackageInfo) error
	Pull(templates.PackageInfo) error
}
