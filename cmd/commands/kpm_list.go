package commands

import (
	"github.com/rohitramu/kpm/cmd/commands/command_groups"
	"github.com/rohitramu/kpm/cmd/flags"
	"github.com/rohitramu/kpm/pkg"
	"github.com/spf13/cobra"
)

var ListCmd = newKpmCommandBuilder(&cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Lists all packages that are currently available for use.",
	GroupID: command_groups.PackageManagement.ID,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Flags
		var kpmHomeDir = flags.KpmHomeDirFlag.GetValue()

		return pkg.ListCmd(kpmHomeDir)
	},
}).Build()
