package flags

import (
	"fmt"

	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
)

var OutputDir = types.NewFlagBuilder[string]("output-dir").
	SetAlias('d').
	SetShortDescription(fmt.Sprintf(
		"Directory in which output files should be written (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<output-name>\" will be deleted if it exists.",
		constants.GeneratedDirName,
	)).
	SetDefaultValueFunc(func(config *config.KpmConfig) string {
		var outputDir, err = files.GetAbsolutePath(constants.GeneratedDirName)
		if err != nil {
			log.Panicf("Failed to get default output directory.")
		}

		return outputDir
	}).
	Build()
