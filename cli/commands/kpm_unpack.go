package commands

import (
	"fmt"

	"github.com/rohitramu/kpm/cli/commands/command_groups"
	"github.com/rohitramu/kpm/cli/flags"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/spf13/cobra"
)

var UnpackCmd = newKpmCommandBuilder(&cobra.Command{
	Use:     "unpack <package name>",
	Short:   "Validates a template package and makes it available for use from the local KPM package repository.",
	GroupID: command_groups.PackageManagement.ID,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Flags
		var skipConfirmation = flags.SkipUserConfirmationFlag.GetValue()
		var kpmHomeDir = flags.KpmHomeDirFlag.GetValue()
		var packageVersion = flags.PackageVersionFlag.GetValue()
		var exportDir = flags.ExportDirFlag.GetValue()
		var exportName = flags.ExportNameFlag.GetValue()

		// Args
		var packageName = args[0]

		// Validation
		{
			// Package version
			if !flags.PackageVersionFlag.IsSetByUser(cmd) {
				// Since the package version was not provided, check the local repository for the highest version
				var err error
				if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
					return fmt.Errorf("package version must be provided if the package does not exist in the local repository: %s", err)
				}
			}

			// Export name
			if !flags.ExportNameFlag.IsSetByUser(cmd) {
				exportName = template_package.GetDefaultExportName(packageName, packageVersion)
			}
		}

		return pkg.UnpackCmd(packageName, packageVersion, exportDir, exportName, kpmHomeDir, skipConfirmation)
	},
}).AddLocalFlags(
	flags.PackageVersionFlag,
	flags.ExportDirFlag,
	flags.ExportNameFlag,
	flags.SkipUserConfirmationFlag,
).Build()
