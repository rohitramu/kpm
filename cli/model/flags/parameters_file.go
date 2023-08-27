package flags

import (
	"github.com/rohitramu/kpm/cli/model/utils/types"
)

var ParametersFile = types.NewFlagBuilder[string]("parameters-file").
	SetAlias('p').
	SetShortDescription("Filepath of the parameters file to use.").
	Build()
