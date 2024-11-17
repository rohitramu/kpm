package cmd_kpm

import (
	"fmt"

	"github.com/rohitramu/kpm/src/cli/model/args"
	"github.com/rohitramu/kpm/src/cli/model/flags"
	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/cli/model/utils/directories"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg"
	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
)

var Inspect = &types.Command{
	Name:             constants.CmdInspect,
	ShortDescription: "Prints the contents of the default parameters file in a template package.",
	Flags: types.FlagCollection{
		BoolFlags: []types.Flag[bool]{
			flags.UserConfirmation,
		},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{args.PackageName("The name of the template package to run.")},
		OptionalArg:   args.PackageVersion("The version of the template package to run."),
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Args
		var packageName = args.MandatoryArgs[0].Value
		var packageVersion = args.OptionalArg.Value

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
			return err
		}

		// Validation
		{
			// Package version
			if packageVersion == "" {
				// Since the package version was not provided, check the local repository for the highest version.
				var err error
				if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
					return fmt.Errorf("could not find package '%s' in the local KPM repository: %s", packageName, err)
				}
			}
		}

		return pkg.InspectCmd(packageName, packageVersion, kpmHomeDir)
	},
}
