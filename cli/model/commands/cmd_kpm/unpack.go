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

var Unpack = &types.Command{
	Name:             constants.CmdUnpack,
	ShortDescription: "Exports a template package to the specified location.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.PackageVersion,
			flags.ExportDir,
			flags.ExportName,
		},
		BoolFlags: []types.Flag[bool]{
			flags.UserConfirmation,
		},
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		// Flags
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)
		var packageVersion = flags.PackageVersion.GetValueOrDefault(config)
		var exportDir = flags.ExportDir.GetValueOrDefault(config)
		var exportName = flags.ExportName.GetValueOrDefault(config)

		// Args
		var packageName = args.MandatoryArgs[0].Value

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
