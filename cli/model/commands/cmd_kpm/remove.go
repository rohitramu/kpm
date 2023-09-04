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

var Remove = &types.Command{
	Name:             constants.CmdRemove,
	Alias:            "rm",
	ShortDescription: "Removes a template package.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{flags.PackageVersion},
		BoolFlags:   []types.Flag[bool]{flags.UserConfirmation},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{{
			Name:             "package-name",
			ShortDescription: "The name of the template package to remove.",
		}},
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)
		var packageVersion = flags.PackageVersion.GetValueOrDefault(config)

		// Args
		var packageName string = args.MandatoryArgs[0].Value

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
			return err
		}

		if packageVersion == "" {
			// Since the package version was not provided, check the local repository for the highest version.
			var err error
			if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
				return fmt.Errorf("could not find package '%s' in the local KPM repository: %s", packageName, err)
			}
		}

		return pkg.RemoveCmd(
			packageName,
			packageVersion,
			kpmHomeDir,
			skipConfirmation)
	},
}
