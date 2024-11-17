package template_repository

import (
	"errors"
	"fmt"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

type repositoryInfo struct {
	Name           string    `yaml:"name"`
	Type           string    `yaml:"type"`
	ConnectionInfo yaml.Node `yaml:"connection"`
}

type RepositoryInfoCollection []*repositoryInfo

type repoInfoParsingFunc func(*repositoryInfo) (Repository, error)

var repoTypeToParsingFunc = map[string]repoInfoParsingFunc{
	repositoryTypeNameFilesystem: repoInfoToFilesystemRepo,
	repositoryTypeNameDocker:     repoInfoToDockerRepo,
}

func (repoInfos RepositoryInfoCollection) ToRepositoryCollection() (*RepositoryCollection, error) {
	var errs []error
	var result = &RepositoryCollection{repos: *linkedhashmap.New()}

	// Create a Repository object based on the user-provided repository information.
	for repoNum, repoInfo := range repoInfos {
		var repo Repository
		var err error
		if repo, err = repoInfo.ToRepository(); err != nil {
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

func (repoInfo *repositoryInfo) ToRepository() (result Repository, err error) {
	// Get the parsing function.
	var parsingFunc, ok = repoTypeToParsingFunc[repoInfo.Type]
	if !ok {
		return nil, fmt.Errorf(
			"unknown repository type '%s' (valid values: %q)",
			repoInfo.Type,
			maps.Keys(repoTypeToParsingFunc),
		)
	}

	// Parse the repo info.
	if result, err = parsingFunc(repoInfo); err != nil {
		return nil, fmt.Errorf(
			"failed to parse information for repository named '%s': %s",
			repoInfo.Name,
			err,
		)
	}

	return result, nil
}
