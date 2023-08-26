package template_repository

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/templates"
)

type Repository interface {
	GetName() string
	GetType() string
	FindPackages(searchTerm string) ([]*templates.PackageInfo, error)
	PackageVersions(packageName string) ([]string, error)
	Push(kpmHomeDir string, packageInfo *templates.PackageInfo) error
	Pull(kpmHomeDir string, packageInfo *templates.PackageInfo) error
}

type PackageNotFoundError struct {
	PackageInfo templates.PackageInfo
}

func (err PackageNotFoundError) Error() string {
	return fmt.Sprintf("package '%s' not found", err.PackageInfo)
}

func (err PackageNotFoundError) Is(target error) bool {
	_, ok := target.(PackageNotFoundError)
	return ok
}
