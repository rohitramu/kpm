package commands

import (
	"fmt"

	"github.com/rohitramu/kpm/cmd/commands/command_groups"
	"github.com/rohitramu/kpm/cmd/flags"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/common"
	"github.com/spf13/cobra"
)

var RunCmd = newKpmCommandBuilder(&cobra.Command{
	Use:     "run <package name>",
	Short:   "Runs a template package.",
	GroupID: command_groups.PackageManagement.ID,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Args
		var packageName = args[0]

		// Flags
		var packageVersion = flags.PackageVersionFlag.GetValue()
		var paramFile = flags.ParametersFileFlag.GetValue()
		var outputDir = flags.OutputDirFlag.GetValue()
		var outputName = flags.OutputNameFlag.GetValue()
		var kpmHomeDir = flags.KpmHomeDirFlag.GetValue()
		var skipConfirmation = flags.SkipUserConfirmationFlag.GetValue()

		// Validation
		var optionalParamFile = &paramFile
		var optionalOutputName = &outputName
		{
			// Package version
			if !flags.PackageVersionFlag.IsSetByUser(cmd) {
				// Since the package version was not provided, check the local repository for the highest version.
				var err error
				if packageVersion, err = common.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
					return fmt.Errorf("could not find package '%s' in the local KPM repository: %s", packageName, err)
				}
			}
			// Parameters file
			if !flags.ParametersFileFlag.IsSetByUser(cmd) {
				optionalParamFile = nil
			}
			// Output name
			if !flags.OutputNameFlag.IsSetByUser(cmd) {
				optionalOutputName = nil
			}
		}

		return pkg.RunCmd(packageName, packageVersion, optionalParamFile, outputDir, optionalOutputName, kpmHomeDir, skipConfirmation)
	},
}).AddLocalFlags(
	flags.PackageVersionFlag,
	flags.ParametersFileFlag,
	flags.OutputDirFlag,
	flags.OutputNameFlag,
	flags.SkipUserConfirmationFlag,
).Build()
