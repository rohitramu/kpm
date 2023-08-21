package template_repository

import (
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

type Repository interface {
	GetName() string
	GetType() string
	Packages() ([]*templates.PackageInfo, error)
	PackageVersions(packageName string) ([]string, error)
	Push(kpmHomeDir string, packageInfo *templates.PackageInfo, userHasConfirmed bool) error
	Pull(kpmHomeDir string, packageInfo *templates.PackageInfo, userHasConfirmed bool) error
}
