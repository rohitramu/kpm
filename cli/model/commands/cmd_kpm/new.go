package cmd_kpm

import (
	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
)

var New = &types.Command{
	Name:             constants.CmdNewPackage,
	Alias:            "new",
	ShortDescription: "Creates a new template package.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.NewPackageOutputDir,
		},
		BoolFlags: []types.Flag[bool]{
			flags.UserConfirmation,
		},
	},
	Args: types.ArgCollection{
		OptionalArg: &types.Arg{
			Name:             "package-name",
			ShortDescription: "The name of the new package.",
			Value:            "hello-kpm",
		},
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) error {
		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)
		var packageDir = flags.NewPackageOutputDir.GetValueOrDefault(config)

		// Args
		var packageName = args.OptionalArg.Value

		return pkg.NewTemplatePackageCmd(packageName, packageDir, skipConfirmation)
	},
}
