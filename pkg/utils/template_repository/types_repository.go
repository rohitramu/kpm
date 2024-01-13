package template_repository

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/template_package"
)

// TODO: Add a repo type for git.

// TODO: Don't use PackageInfo type for repository operations.

type Repository interface {
	GetName() string
	GetType() string
	FindPackages(ch chan<- *template_package.PackageInfo, searchTerm string) error
	PackageVersions(ch chan<- string, packageName string) error
	Push(kpmHomeDir string, packageInfo *template_package.PackageInfo) error
	Pull(kpmHomeDir string, packageInfo *template_package.PackageInfo) error
}

type PackageNotFoundError struct {
	PackageInfo template_package.PackageInfo
}

func (err PackageNotFoundError) Error() string {
	return fmt.Sprintf("package '%s' not found", err.PackageInfo)
}

func (err PackageNotFoundError) Is(target error) bool {
	_, ok := target.(PackageNotFoundError)
	return ok
}
