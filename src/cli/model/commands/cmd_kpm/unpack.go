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

var Unpack = &types.Command{
	Name:             constants.CmdUnpack,
	ShortDescription: "Exports a template package to the specified location.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.ExportDir,
			flags.ExportName,
		},
		BoolFlags: []types.Flag[bool]{
			flags.UserConfirmation,
		},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{args.PackageName("The name of the template package to unpack.")},
		OptionalArg:   args.PackageVersion("The version of the template package to unpack.  If not provided, the latest package version will be used."),
	},
	ExecuteFunc: func(config *config.KpmConfig, inputArgs types.ArgCollection) (err error) {
		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)
		var exportDir = flags.ExportDir.GetValueOrDefault(config)
		var exportName = flags.ExportName.GetValueOrDefault(config)

		// Args
		var packageName = inputArgs.MandatoryArgs[0].Value
		var packageVersion = inputArgs.OptionalArg.Value

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
			return err
		}

		// Validation
		{
			// Package version
			if packageVersion == "" {
				// Since the package version was not provided, check the local repository for the highest version
				var err error
				if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
					return fmt.Errorf("package version must be provided if the package does not exist in the local repository: %s", err)
				}
			}

			// Export name
			if exportName == "" {
				exportName = template_package.GetDefaultExportName(packageName, packageVersion)
			}
		}

		return pkg.UnpackCmd(packageName, packageVersion, exportDir, exportName, kpmHomeDir, skipConfirmation)
	},
}
