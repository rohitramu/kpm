package flags

import (
	"github.com/rohitramu/kpm/cli/model/utils/types"
)

// TODO: Use optional arg for package version instead of a flag.

var PackageVersion = types.NewFlagBuilder[string]("version").
	SetAlias('v').
	SetShortDescription("The template package's version.").
	SetValidationFunc(ValidatePackageVersion()).
	Build()
