package template_repository

import (
	"errors"
	"fmt"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"gopkg.in/yaml.v3"
)

var _ yaml.Unmarshaler = &RepositoryCollection{}
var _ yaml.Marshaler = &RepositoryCollection{}

type RepositoryCollection struct {
	repos linkedhashmap.Map
}

func (result *RepositoryCollection) UnmarshalYAML(unmarshaller *yaml.Node) (err error) {
	var repoInfos RepositoryInfoCollection

	err = unmarshaller.Decode(&repoInfos)
	if err != nil {
		return err
	}

	var tmp *RepositoryCollection
	tmp, err = repoInfos.ToRepositoryCollection()
	if err != nil {
		return fmt.Errorf("failed to parse repository information: %s", err)
	}

	*result = *tmp

	return nil
}

func (rc *RepositoryCollection) MarshalYAML() (any, error) {
	return nil, errors.New("not yet implemented")
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
	// Get the repo from the map.
	var repoUncasted, found = rc.repos.Get(repoName)
	if !found {
		return nil, fmt.Errorf("unknown repository '%s'", repoName)
	}

	// Cast it to the correct type.
	var repo = castObjToRepo(repoUncasted)

	return repo, nil
}

func (rc *RepositoryCollection) AddRepository(repo Repository) error {
	var repoName = repo.GetName()
	if _, found := rc.repos.Get(repoName); found {
		return fmt.Errorf("repository named '%s' already exists", repoName)
	}

	rc.repos.Put(repoName, repo)
	return nil
}

func (rc *RepositoryCollection) AddRepositoryInfo(repoInfo repositoryInfo) (err error) {
	var repo Repository

	repo, err = repoInfo.ToRepository()
	if err != nil {
		return err
	}

	return rc.AddRepository(repo)
}

func (rc *RepositoryCollection) RemoveRepository(repoName string) error {
	if _, found := rc.repos.Get(repoName); !found {
		return fmt.Errorf("unknown repository '%s'", repoName)
	}

	rc.repos.Remove(repoName)
	return nil
}

func castObjToRepo(uncasted any) Repository {
	var repo, ok = uncasted.(Repository)
	if !ok {
		log.Panicf("Failed to cast the repo to a Repository type.")
	}

	return repo
}
