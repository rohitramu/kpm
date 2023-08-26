package template_repository

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/templates"
)

type ErrPackageNotFoundType struct {
	PackageInfo templates.PackageInfo
}

var ErrPackageNotFound = ErrPackageNotFoundType{}

func (err ErrPackageNotFoundType) Error() string {
	return fmt.Sprintf("failed to find package '%s'", err.PackageInfo)
}

func (thisErr ErrPackageNotFoundType) Is(target error) bool {
	var _, ok = target.(ErrPackageNotFoundType)

	// Don't worry about which package we couldn't find - we only want to compare error types here.
	return ok
}
