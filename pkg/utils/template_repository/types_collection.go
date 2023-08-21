package template_repository

import (
	"errors"
	"fmt"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

type RepositoryCollection interface {
	RepositoryNames() []string
	Pull(templates.PackageInfo) error
	Push(repositoryName string, packageInfo templates.PackageInfo) error
}

type repositoryCollection struct {
	repos linkedhashmap.Map
}

func (rc *repositoryCollection) RepositoryNames() []string {
	var result = make([]string, 0, rc.repos.Size())
	for _, repoNameUncasted := range rc.repos.Keys() {
		var repoName, ok = repoNameUncasted.(string)
		if !ok {
			log.Panicf("Failed to cast repository name to string.")
		}

		result = append(result, repoName)
	}

	return result
}

func (rc *repositoryCollection) Pull(packageInfo templates.PackageInfo) error {
	var err error

	var it = rc.repos.Iterator()
	it.Begin()
	for it.Next() {
		// Get repo.
		var repo = castObjToRepo(it.Value())

		// Attempt to pull package from repo.
		log.Verbosef(
			"Pulling package '%s', version '%s' from repository '%s'.",
			packageInfo.Name,
			packageInfo.Version, repo.GetName(),
		)
		err = repo.Pull(packageInfo)
		if err == nil {
			// We succesfully pulled the package
			return nil
		}

		// Failed to pull package from this repo - keep checking other repos.
		if errors.Is(err, ErrPackageNotFound) {
			log.Infof(
				"Could not find package '%s', version '%s' in repository '%s'.",
				packageInfo.Name,
				packageInfo.Version,
				repo.GetName(),
			)
		}
	}

	// If we get to this point, that means we didn't find the package in any repos.
	return fmt.Errorf(
		"failed to find package in any repositories: %s",
		ErrPackageNotFoundType{PackageInfo: packageInfo},
	)
}

func (rc *repositoryCollection) Push(repoName string, packageInfo templates.PackageInfo) error {
	// Get the repo from the map.
	var repo, err = rc.getRepo(repoName)
	if err != nil {
		return err
	}

	// Push the template package to the repository.
	err = repo.Push(packageInfo)

	return err
}

func (rc *repositoryCollection) getRepo(repoName string) (Repository, error) {
	// Get the repo from the map.
	var repoUncasted, found = rc.repos.Get(repoName)
	if !found {
		return nil, fmt.Errorf("unknown repository '%s'", repoName)
	}

	// Cast it to the correct type.
	var repo = castObjToRepo(repoUncasted)

	return repo, nil
}

func castObjToRepo(uncasted any) Repository {
	var repo, ok = uncasted.(Repository)
	if !ok {
		log.Panicf("Failed to cast the repo to a Repository type.")
	}

	return repo
}
