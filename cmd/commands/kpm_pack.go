package commands

import (
	"fmt"

	"github.com/rohitramu/kpm/cmd/commands/command_groups"
	"github.com/rohitramu/kpm/cmd/flags"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/spf13/cobra"
)

var PackCmd = newKpmCommandBuilder(&cobra.Command{
	Use:     "pack [<package directory>]",
	Short:   "Validates a template package and makes it available for use from the local KPM package repository.",
	GroupID: command_groups.PackageManagement.ID,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Flags
		var skipConfirmation = flags.SkipUserConfirmationFlag.GetValue()
		var kpmHomeDir = flags.KpmHomeDirFlag.GetValue()

		// Args
		var packageDir string
		if len(args) > 0 {
			packageDir = args[0]
		} else {
			var err error
			packageDir, err = files.GetWorkingDir()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %s", err)
			}
		}

		return pkg.PackCmd(packageDir, kpmHomeDir, skipConfirmation)
	},
}).AddLocalFlags(
	flags.SkipUserConfirmationFlag,
).Build()
