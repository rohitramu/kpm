package cmd_kpm

import (
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

var List = &types.Command{
	Name:             constants.CmdList,
	Alias:            "ls",
	ShortDescription: "Lists all template packages.",
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		// Get the list of package names.
		var packages []string
		packages, err = pkg.ListCmd(kpmHomeDir)
		if err != nil {
			return err
		}

		// Print package names.
		for _, packageName := range packages {
			log.Outputf(packageName)
		}

		return nil
	},
}
