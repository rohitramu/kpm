package flags

import (
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
)

var ExportName = types.NewFlagBuilder[string]("export-name").
	SetAlias('n').
	SetShortDescription("Name of the exported template package (defaults to \"<package name>-<package version>\").").
	Build()
