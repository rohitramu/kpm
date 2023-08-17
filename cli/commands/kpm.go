package commands

import (
	"github.com/spf13/cobra"

	"github.com/rohitramu/kpm/cli/commands/command_groups"
	"github.com/rohitramu/kpm/cli/flags"
)

var versionString = "0.0.0"

// KpmCmd is the entrypoint for this application.
var KpmCmd = newKpmCommandBuilder(&cobra.Command{
	Use:     "kpm",
	Short:   "KPM is a text generation tool.",
	Version: versionString,
}).AddCommandGroups(
	command_groups.PackageManagement,
	command_groups.PackageEditing,
).AddPersistentFlags(
	flags.LogLevelFlag,
).AddSubcommands(
	ListCmd,
	RemoveCmd,
	PurgeCmd,
	PackCmd,
	UnpackCmd,
	InspectCmd,
	RunCmd,
	NewCmd,
).Build()
