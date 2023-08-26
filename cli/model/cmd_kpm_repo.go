package model

import (
	"errors"
	"fmt"

	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/directories"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_repository"
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

var repoRemotesCmd = &Command{
	Name:             constants.CmdRepoList,
	Alias:            "ls",
	ShortDescription: "Lists the names of available repositories.",
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var repos []string
		repos, err = pkg.ListPackageRepositories(&config.Repositories)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			log.Outputf(repo)
		}

		return nil
	},
}

var repoFindCmd = &Command{
	Name:             constants.CmdRepoFind,
	ShortDescription: "Finds packages.",
	Flags: FlagCollection{StringFlags: []Flag[string]{
		repoNameFlag,
	}},
	Args: ArgCollection{OptionalArg: &Arg{
		Name:             "search-term",
		ShortDescription: "A search term to use for finding packages.",
	}},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var searchTerm = ""
		if args.OptionalArg != nil {
			searchTerm = args.OptionalArg.Value
		}

		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		var packageInfos []*templates.PackageInfo
		packageInfos, err = pkg.FindPackages(
			kpmHomeDir,
			&config.Repositories,
			repoNameFlag.GetValueOrDefault(config),
			searchTerm,
		)
		if err != nil {
			return err
		}

		for _, packageInfo := range packageInfos {
			log.Outputf("%s", packageInfo)
		}

		return nil
	},
}

var repoPushCmd = &Command{
	Name:             constants.CmdRepoPush,
	ShortDescription: "Pushes a template package to a repository.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			repoNameFlag,
			packageVersionFlag,
		},
		BoolFlags: []Flag[bool]{
			userConfirmationFlag,
		},
	},
	Args: ArgCollection{MandatoryArgs: []*Arg{
		{
			Name:             "package-name",
			ShortDescription: "The name of the template package to push.",
			IsValidFunc:      validatePackageName,
		},
	}},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		var packageName = args.MandatoryArgs[0].Value
		var packageVersion = packageVersionFlag.GetValueOrDefault(config)

		// If the repo name isn't provided, pick the first one.
		var repoName = repoNameFlag.GetValueOrDefault(config)
		if repoName == "" {
			var repoNames = config.Repositories.GetRepositoryNames()
			if len(repoNames) == 0 {
				return errors.New("no repositories configured")
			}

			repoName = repoNames[0]
		}

		return pkg.PushPackage(
			kpmHomeDir,
			&config.Repositories,
			repoName,
			packageName,
			packageVersion,
			userConfirmationFlag.GetValueOrDefault(config),
		)
	},
}

var repoPullCmd = &Command{
	Name:             constants.CmdRepoPull,
	ShortDescription: "Pulls a template package from a repository.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			repoNameFlag,
			packageVersionFlag,
		},
	},
	Args: ArgCollection{MandatoryArgs: []*Arg{
		{
			Name:             "package-name",
			ShortDescription: "The name of the template package to pull.",
			IsValidFunc:      validatePackageName,
		},
	}},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var repoNames = config.Repositories.GetRepositoryNames()
		if len(repoNames) == 0 {
			return errors.New("no repositories configured")
		}

		var kpmHomeDir string
		if kpmHomeDir, err = directories.GetKpmHomeDir(); err != nil {
			return err
		}

		var packageName = args.MandatoryArgs[0].Value
		var packageVersion = packageVersionFlag.GetValueOrDefault(config)

		// If the repo name isn't provided, pick the first one.
		var repoName = repoNameFlag.GetValueOrDefault(config)
		if repoName != "" {
			return pkg.PullPackage(
				kpmHomeDir,
				&config.Repositories,
				repoName,
				packageName,
				packageVersion,
			)
		}

		for _, repoName := range repoNames {
			err = pkg.PullPackage(
				kpmHomeDir,
				&config.Repositories,
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
