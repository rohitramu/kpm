package template_repository

import (
	"errors"
	"fmt"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"golang.org/x/exp/maps"
)

type repoInfoParsingFunc func(RepositoryInfo) (Repository, error)

var repoTypeToParsingFunc = map[string]repoInfoParsingFunc{
	repositoryTypeNameFilesystem: repoInfoToFilesystemRepo,
	repositoryTypeNameDocker:     repoInfoToDockerRepo,
}

func GetRepositoriesFromInfo(repoInfos ...RepositoryInfo) (RepositoryCollection, error) {
	var errs []error
	var result = &repositoryCollection{repos: *linkedhashmap.New()}

	// Create a Repository object based on the user-provided repository information.
	for repoNum, repoInfo := range repoInfos {
		var repo Repository
		var err error
		if repo, err = parseRepoInfo(result, repoInfo); err != nil {
			// Collect the error and continue parsing (maybe the caller is happy to ignore bad repo references).
			errs = append(errs, fmt.Errorf("#%d: %s", repoNum, err))
			continue
		}

		result.repos.Put(repoInfo.Name, repo)
	}

	// Combine the list of error messages if there were any.
	var err error
	if len(errs) > 0 {
		var combinedErrs = errors.Join(errs...)
		err = errors.Join(errors.New("failed to parse information for some repositories"), combinedErrs)
	}

	return result, err
}

func parseRepoInfo(repoCollection *repositoryCollection, repoInfo RepositoryInfo) (Repository, error) {
	var err error
	var result Repository

	// Validate repo name.
	if err = validateRepoName(repoInfo.Name); err != nil {
		return result, fmt.Errorf("invalid repository name: %s", err)
	}

	// Make sure a repo with the same name doesn't already exist.
	if err = ensureRepoHasntBeenParsed(repoCollection, repoInfo.Name); err != nil {
		return result, err
	}

	var parsingFunc, ok = repoTypeToParsingFunc[repoInfo.Type]
	if !ok {
		return result, fmt.Errorf("unknown repository type '%s' (valid values: %q)", repoInfo.Type, maps.Keys(repoTypeToParsingFunc))
	}

	if result, err = parsingFunc(repoInfo); err != nil {
		return result, fmt.Errorf("failed to parse information for repository named '%s': %s", repoInfo.Name, err)
	}

	return result, nil
}

func validateRepoName(repoName string) error {
	if len(repoName) == 0 {
		return fmt.Errorf("repository names cannot be empty")
	}

	return nil
}

func ensureRepoHasntBeenParsed(repoCollection *repositoryCollection, repoName string) error {
	if existingRepoUncasted, alreadyExists := repoCollection.repos.Get(repoName); alreadyExists {
		// Found a repo that already exists.  Cast it so we can get its type name.
		if existingRepo, ok := existingRepoUncasted.(Repository); ok {
			return fmt.Errorf("a '%s' repository with name '%s' already exists", existingRepo.GetType(), repoName)
		} else {
			return fmt.Errorf("repository with name '%s' already exists", repoName)
		}
	}

	return nil
}
