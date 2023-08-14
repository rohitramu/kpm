package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rohitramu/kpm/cmd/commands/command_groups"
	"github.com/rohitramu/kpm/cmd/flags"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/files"
)

var NewCmd = newKpmCommandBuilder(&cobra.Command{
	Use:     "new-package <package name> [<package directory>]",
	Aliases: []string{"new"},
	Short:   "Creates a new template package.",
	GroupID: command_groups.PackageEditing.ID,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Flags
		var skipConfirmation = flags.SkipUserConfirmationFlag.GetValue()

		// Args
		var packageName = args[0]
		var packageDir string
		if len(args) > 1 {
			packageDir = args[1]
		} else {
			var err error
			packageDir, err = files.GetWorkingDir()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %s", err)
			}
		}

		return pkg.NewTemplatePackageCmd(packageName, packageDir, skipConfirmation)
	},
}).AddLocalFlags(
	flags.SkipUserConfirmationFlag,
).Build()
