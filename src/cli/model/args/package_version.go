package args

import (
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
)

func PackageVersion(shortDescription string) *types.Arg {
	return &types.Arg{
		Name:             "package-version",
		ShortDescription: shortDescription,
		Value:            "",
		IsValidFunc:      validation.ValidatePackageVersion,
	}
}
