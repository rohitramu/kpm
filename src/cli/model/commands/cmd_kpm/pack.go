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

var Pack = &types.Command{
	Name:             constants.CmdPack,
	ShortDescription: "Validates a template package and makes it available for use.",
	Flags: types.FlagCollection{
		BoolFlags: []types.Flag[bool]{flags.UserConfirmation},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{args.PackageDirectory("The location of the template package directory which should be packed.")},
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Args
		var packageDir = args.MandatoryArgs[0].Value

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
			return err
		}

		return pkg.PackCmd(packageDir, kpmHomeDir, skipConfirmation)
	},
}
