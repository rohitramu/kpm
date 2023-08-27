package cmd_kpm_repo

import (
	"errors"
	"fmt"

	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/template_repository"
	"github.com/rohitramu/kpm/pkg/utils/validation"
)

var PullCmd = &types.Command{
	Name:             constants.CmdRepoPull,
	ShortDescription: "Pulls a template package from a repository.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.RepoName,
			flags.PackageVersion,
		},
	},
	Args: types.ArgCollection{MandatoryArgs: []*types.Arg{
		{
			Name:             "package-name",
			ShortDescription: "The name of the template package to pull.",
			IsValidFunc:      validation.ValidatePackageName,
		},
	}},
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		var repoNames = config.Repositories.GetRepositoryNames()
		if len(repoNames) == 0 {
			return errors.New("no repositories configured")
		}

		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		var packageName = args.MandatoryArgs[0].Value
		var packageVersion = flags.PackageVersion.GetValueOrDefault(config)

		// If the repo name isn't provided, pick the first one.
		var repoName = flags.RepoName.GetValueOrDefault(config)
		if repoName != "" {
			return pkg.PullPackage(
				kpmHomeDir,
				config.Repositories,
				repoName,
				packageName,
				packageVersion,
			)
		}

		for _, repoName := range repoNames {
			err = pkg.PullPackage(
				kpmHomeDir,
				config.Repositories,
				repoName,
				packageName,
				packageVersion,
			)

			if err == nil {
				// We found and pulled the package.
				return nil
			}

			if errors.Is(err, template_repository.PackageNotFoundError{}) {
				// If the package wasn't found, continue checking other repos.
				continue
			} else {
				// Something failed in the process - return the error.
				return err
			}
		}

		return fmt.Errorf(
			"cannot find package '%s' (version '%s') in configured repsitories",
			packageName,
			packageVersion,
		)
	},
}
