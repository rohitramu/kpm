package cmd_kpm_repo

import (
	"errors"

	"github.com/rohitramu/kpm/src/cli/model/args"
	"github.com/rohitramu/kpm/src/cli/model/flags"
	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/cli/model/utils/directories"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg"
)

var PushCmd = &types.Command{
	Name:             constants.CmdRepoPush,
	ShortDescription: "Pushes a template package to a repository.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.RepoName,
		},
		BoolFlags: []types.Flag[bool]{
			flags.UserConfirmation,
		},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{
			args.PackageName("The name of the template package to push."),
			args.PackageVersion("The version of the template package to push."),
		},
	},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		// Flags
		var repoName = flags.RepoName.GetValueOrDefault(config)
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Args
		var packageName = args.MandatoryArgs[0].Value
		var packageVersion = args.MandatoryArgs[1].Value

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
			return err
		}

		// Validation
		{
			// If the repo name isn't provided, pick the first one.
			if repoName == "" {
				var repoNames = config.Repositories.GetRepositoryNames()
				if len(repoNames) == 0 {
					return errors.New("no repositories configured")
				}

				repoName = repoNames[0]
			}
		}

		return pkg.PushPackage(
			kpmHomeDir,
			config.Repositories,
			repoName,
			packageName,
			packageVersion,
			skipConfirmation,
		)
	},
}
