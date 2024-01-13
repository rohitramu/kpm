package args

import (
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg/utils/validation"
)

func PackageDirectory(shortDescription string) *types.Arg {
	return &types.Arg{
		Name:             "package-directory",
		ShortDescription: shortDescription,
		Value:            "",
		IsValidFunc:      validation.ValidatePackageName,
	}
}
