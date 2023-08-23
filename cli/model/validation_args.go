package model

import (
	"github.com/rohitramu/kpm/pkg/utils/validation"
)

var validatePackageName = func(packageName string) error {
	return validation.ValidatePackageName(packageName)
}
