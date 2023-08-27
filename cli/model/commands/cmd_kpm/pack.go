package cmd_kpm

import (
	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
)

var Pack = &types.Command{
	Name:             constants.CmdPack,
	ShortDescription: "Validates a template package and makes it available for use.",
	Flags: types.FlagCollection{
		BoolFlags: []types.Flag[bool]{flags.UserConfirmation},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{{
			Name:             "package-directory",
			ShortDescription: "The location of the template package directory which should be packed.",
		}},
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Args
		var packageDir = args.MandatoryArgs[0].Value

		return pkg.PackCmd(packageDir, kpmHomeDir, skipConfirmation)
	},
}
