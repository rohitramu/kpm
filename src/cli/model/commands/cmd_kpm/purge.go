package cmd_kpm

import (
	"github.com/rohitramu/kpm/src/cli/model/args"
	"github.com/rohitramu/kpm/src/cli/model/flags"
	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/cli/model/utils/directories"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg"
)

var Purge = &types.Command{
	Name:             constants.CmdPurge,
	ShortDescription: "Removes all versions of a template package.",
	Flags: types.FlagCollection{
		BoolFlags: []types.Flag[bool]{flags.UserConfirmation},
	},
	Args: types.ArgCollection{
		OptionalArg: args.PackageName("The name of the template package to purge.  If this is not provided, all versions of all template packages will be deleted."),
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Args
		var packageName string
		if args.OptionalArg != nil {
			packageName = args.OptionalArg.Value
		}

		return pkg.PurgeCmd(packageName, skipConfirmation, kpmHomeDir)
	},
}
