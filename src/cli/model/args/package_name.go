package args

import (
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
)

func PackageName(shortDescription string) *types.Arg {
	return PackageNameWithDefault(shortDescription, "")
}

func PackageNameWithDefault(shortDescription string, defaultValue string) *types.Arg {
	return &types.Arg{
		Name:             "package-name",
		ShortDescription: shortDescription,
		Value:            defaultValue,
		IsValidFunc:      validation.ValidatePackageName,
	}
}
