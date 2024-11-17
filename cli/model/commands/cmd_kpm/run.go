package cmd_kpm

import (
	"fmt"

	"github.com/rohitramu/kpm/cli/model/args"
	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
)

var Run = &types.Command{
	Name:             constants.CmdRun,
	ShortDescription: "Runs a template package.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.ParametersFile,
			flags.OutputDir,
			flags.OutputName,
		},
		BoolFlags: []types.Flag[bool]{
			flags.UserConfirmation,
		},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{args.PackageName("The name of the template package to run.")},
		OptionalArg:   args.PackageVersion("The version of the template package to run.  If not set, the latest version will be run."),
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		// Args
		var packageName = args.MandatoryArgs[0].Value
		var packageVersion = args.OptionalArg.Value

		// Flags
		var paramFile = flags.ParametersFile.GetValueOrDefault(config)
		var outputDir = flags.OutputDir.GetValueOrDefault(config)
		var outputName = flags.OutputName.GetValueOrDefault(config)
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
			return err
		}

		// Validation
		var optionalParamFile = &paramFile
		var optionalOutputName = &outputName
		{
			// Package version
			if packageVersion == "" {
				// Since the package version was not provided, check the local repository for the highest version.
				var err error
				if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
					return fmt.Errorf("could not find package '%s' in the local KPM repository: %s", packageName, err)
				}
			}
			// Parameters file
			if paramFile == "" {
				optionalParamFile = nil
			}
			// Output name
			if outputName == "" {
				optionalOutputName = nil
			}
		}

		return pkg.RunCmd(packageName, packageVersion, optionalParamFile, outputDir, optionalOutputName, kpmHomeDir, skipConfirmation)
	},
}
