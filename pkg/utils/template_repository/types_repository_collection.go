package template_repository

import (
	"fmt"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"gopkg.in/yaml.v3"
)

type RepositoryCollection struct {
	repos linkedhashmap.Map
}

func (result *RepositoryCollection) UnmarshalYAML(unmarshaller *yaml.Node) (err error) {
	var repoInfos *[]*RepositoryInfo

	err = unmarshaller.Decode(&repoInfos)
	if err != nil {
		return err
	}

	var tmp *RepositoryCollection
	tmp, err = GetRepositoriesFromInfo(*repoInfos...)
	if err != nil {
		return fmt.Errorf("failed to parse repository information: %s", err)
	}

	*result = *tmp

	return nil
}

func (rc *RepositoryCollection) GetRepositoryNames() []string {
	var result = make([]string, 0, rc.repos.Size())
	var it = rc.repos.Iterator()
	for it.Next() {
		var repoName, ok = it.Key().(string)
		if !ok {
			log.Panicf("Failed to cast repository name to string.")
		}

		result = append(result, repoName)
	}

	return result
}

func (rc *RepositoryCollection) GetRepository(repoName string) (Repository, error) {
	var repo, err = rc.getRepo(repoName)
	return repo, err
}

func (rc *RepositoryCollection) getRepo(repoName string) (Repository, error) {
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
