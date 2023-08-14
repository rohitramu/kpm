package commands

import (
	"fmt"

	"github.com/rohitramu/kpm/cmd/commands/command_groups"
	"github.com/rohitramu/kpm/cmd/flags"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/common"
	"github.com/spf13/cobra"
)

var RemoveCmd = newKpmCommandBuilder(&cobra.Command{
	Use:     "remove <package name>",
	Aliases: []string{"rm"},
	Short:   "Removes a package from the local KPM package repository.",
	GroupID: command_groups.PackageManagement.ID,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Flags
		var kpmHomeDir = flags.KpmHomeDirFlag.GetValue()
		var shouldSkipUserConfirmation = flags.SkipUserConfirmationFlag.GetValue()
		var packageVersion = flags.PackageVersionFlag.GetValue()

		// Args
		var packageName string = args[0]

		if !flags.PackageVersionFlag.IsSetByUser(cmd) {
			// Since the package version was not provided, check the local repository for the highest version.
			var err error
			if packageVersion, err = common.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
				return fmt.Errorf("could not find package '%s' in the local KPM repository: %s", packageName, err)
			}
		}

		return pkg.RemoveCmd(
			packageName,
			packageVersion,
			kpmHomeDir,
			shouldSkipUserConfirmation)
	},
}).AddLocalFlags(
	flags.PackageVersionFlag,
	flags.SkipUserConfirmationFlag,
).Build()
