package cmd_kpm_repo

import (
	"errors"

	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/validation"
)

var PushCmd = &types.Command{
	Name:             constants.CmdRepoPush,
	ShortDescription: "Pushes a template package to a repository.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.RepoName,
			flags.PackageVersion,
		},
		BoolFlags: []types.Flag[bool]{
			flags.UserConfirmation,
		},
	},
	Args: types.ArgCollection{MandatoryArgs: []*types.Arg{
		{
			Name:             "package-name",
			ShortDescription: "The name of the template package to push.",
			IsValidFunc:      validation.ValidatePackageName,
		},
	}},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		var packageName = args.MandatoryArgs[0].Value
		var packageVersion = flags.PackageVersion.GetValueOrDefault(config)

		// If the repo name isn't provided, pick the first one.
		var repoName = flags.RepoName.GetValueOrDefault(config)
		if repoName == "" {
			var repoNames = config.Repositories.GetRepositoryNames()
			if len(repoNames) == 0 {
				return errors.New("no repositories configured")
			}

			repoName = repoNames[0]
		}

		return pkg.PushPackage(
			kpmHomeDir,
			config.Repositories,
			repoName,
			packageName,
			packageVersion,
			flags.UserConfirmation.GetValueOrDefault(config),
		)
	},
}
