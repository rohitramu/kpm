package cmd_kpm

import (
	"github.com/rohitramu/kpm/src/cli/model/args"
	"github.com/rohitramu/kpm/src/cli/model/flags"
	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg"
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
		OptionalArg: args.PackageNameWithDefault("The name of the new package.", "hello-kpm"),
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
