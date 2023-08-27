package flags

import (
	"github.com/rohitramu/kpm/cli/model/utils/types"
)

var OutputName = types.NewFlagBuilder[string]("output-name").
	SetAlias('n').
	SetShortDescription("Name of the output (defaults to \"<package name>-<package version>\").").
	Build()
