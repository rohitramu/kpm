package cmd_kpm

import (
	"github.com/rohitramu/kpm/src/cli/model/flags"
	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/cli/model/utils/directories"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
)

var List = &types.Command{
	Name:             constants.CmdList,
	Alias:            "ls",
	ShortDescription: "Lists all template packages.",
	Flags: types.FlagCollection{
		BoolFlags: []types.Flag[bool]{flags.UserConfirmation},
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
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
