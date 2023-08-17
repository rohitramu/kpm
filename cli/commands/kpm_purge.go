package commands

import (
	"github.com/rohitramu/kpm/cli/commands/command_groups"
	"github.com/rohitramu/kpm/cli/flags"
	"github.com/rohitramu/kpm/pkg"
	"github.com/spf13/cobra"
)

var PurgeCmd = newKpmCommandBuilder(&cobra.Command{
	Use:     "purge [<package directory>]",
	Short:   "Removes all versions of a package from the local KPM package repository.",
	GroupID: command_groups.PackageManagement.ID,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Flags
		var skipConfirmation = flags.SkipUserConfirmationFlag.GetValue()
		var kpmHomeDir = flags.KpmHomeDirFlag.GetValue()

		// Args
		var packageName = ""
		if len(args) > 0 {
			packageName = args[0]
		}

		return pkg.PurgeCmd(packageName, skipConfirmation, kpmHomeDir)
	},
}).AddLocalFlags(
	flags.SkipUserConfirmationFlag,
).Build()
