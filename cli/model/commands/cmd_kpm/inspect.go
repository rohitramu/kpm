package cmd_kpm

import (
	"fmt"

	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
)

var Inspect = &types.Command{
	Name:             constants.CmdInspect,
	ShortDescription: "Prints the contents of the default parameters file in a template package.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.PackageVersion,
		},
		BoolFlags: []types.Flag[bool]{
			flags.UserConfirmation,
		},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{{
			Name:             "package-name",
			ShortDescription: "The name of the template package to run.",
		}},
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		// Args
		var packageName = args.MandatoryArgs[0].Value

		// Flags
		var packageVersion = flags.PackageVersion.GetValueOrDefault(config)

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
