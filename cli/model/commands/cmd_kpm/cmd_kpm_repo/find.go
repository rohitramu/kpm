package cmd_kpm_repo

import (
	"github.com/rohitramu/kpm/cli/model/args"
	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
)

var FindCmd = &types.Command{
	Name:             constants.CmdRepoFind,
	ShortDescription: "Finds packages.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{flags.RepoName},
		BoolFlags:   []types.Flag[bool]{flags.UserConfirmation},
	},
	Args: types.ArgCollection{OptionalArg: args.SearchTerm("A search term to use for finding packages.")},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		// Flags
		var repoName = flags.RepoName.GetValueOrDefault(config)
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Args
		var searchTerm = ""
		if args.OptionalArg != nil {
			searchTerm = args.OptionalArg.Value
		}

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
			return err
		}

		// Create a channel for the results.
		ch := make(chan *template_package.PackageInfo, 1)

		// Set up the receiver
		go func() {
			// Print the results.
			numTags := 0
			for packageInfo := range ch {
				numTags++
				log.Outputf(packageInfo.String())
			}
		}()

		err = pkg.FindPackages(
			ch,
			kpmHomeDir,
			config.Repositories,
			repoName,
			searchTerm,
		)
		if err != nil {
			return err
		}

		return nil
	},
}
