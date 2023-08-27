package flags

import (
	"fmt"

	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

var NewPackageOutputDir = types.NewFlagBuilder[string]("output-dir").
	SetAlias('d').
	SetShortDescription(fmt.Sprintf(
		"Directory in which the new template package should be generated (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<output-name>\" will be deleted if it exists.",
		constants.NewTemplatePackageDirName,
	)).
	SetDefaultValueFunc(func(config *config.KpmConfig) string {
		var outputDir, err = files.GetAbsolutePath(constants.NewTemplatePackageDirName)
		if err != nil {
			log.Panicf("Failed to get default output directory.")
		}

		return outputDir
	}).
	Build()
