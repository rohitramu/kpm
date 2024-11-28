package flags

import (
	"fmt"

	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
)

var ExportDir = types.NewFlagBuilder[string]("export-dir").
	SetAlias('d').
	SetShortDescription(fmt.Sprintf(
		"The directory which the template package should be exported to (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<export-name>\" will be deleted if it exists.",
		constants.ExportDirName,
	)).
	SetDefaultValueFunc(func(kc *config.KpmConfig) string { return constants.ExportDirName }).
	SetValidationFunc(ValidateDirExists()).
	Build()
