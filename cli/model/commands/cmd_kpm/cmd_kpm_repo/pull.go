package cmd_kpm_repo

import (
	"errors"
	"fmt"

	"github.com/rohitramu/kpm/cli/model/args"
	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/template_repository"
)

var PullCmd = &types.Command{
	Name:             constants.CmdRepoPull,
	ShortDescription: "Pulls a template package from a repository.",
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.RepoName,
		},
		BoolFlags: []types.Flag[bool]{flags.UserConfirmation},
	},
	Args: types.ArgCollection{
		MandatoryArgs: []*types.Arg{args.PackageName("The name of the template package to pull.")},
		OptionalArg:   args.PackageVersion("The version of the template package to pull."),
	},
	ExecuteFunc: func(config *config.KpmConfig, inputArgs types.ArgCollection) (err error) {
		var repoNames = config.Repositories.GetRepositoryNames()
		if len(repoNames) == 0 {
			return errors.New("no repositories configured")
		}

		// Flags
		var repoName = flags.RepoName.GetValueOrDefault(config)
		var skipConfirmation = flags.UserConfirmation.GetValueOrDefault(config)

		// Args
		var packageName = inputArgs.MandatoryArgs[0].Value
		var packageVersion = inputArgs.OptionalArg.Value

		// Get KPM home directory or create it if it doesn't exist.
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetOrCreateKpmHomeDir(skipConfirmation); err != nil {
			return err
		}

		// If the repo name is provided, don't look in other repos.
		if repoName != "" {
			return pkg.PullPackage(
				kpmHomeDir,
				config.Repositories,
				repoName,
				packageName,
				packageVersion,
			)
		}

		// If the repo name is not provided, search all repos in order of priority.
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
