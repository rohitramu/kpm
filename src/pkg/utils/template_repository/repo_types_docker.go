package template_repository

import (
	"fmt"

	"github.com/rohitramu/kpm/src/pkg/utils/docker"
	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
)

// TODO: Implement Docker repository support.

const repositoryTypeNameDocker = "docker"

var _ Repository = &dockerRepository{}

type dockerRepository struct {
	name           string
	connectionInfo dockerRepositoryConnectionInfo
}

type dockerRepositoryConnectionInfo struct {
	username     string
	organization string
	registry     string
}

func (repo *dockerRepository) GetName() string {
	return repo.name
}

func (repo *dockerRepository) GetType() string {
	return repositoryTypeNameDocker
}

func (repo *dockerRepository) FindPackages(
	ch chan<- *template_package.PackageInfo,
	searchTeam string,
) error {
	return fmt.Errorf("not yet implemented")
}

func (repo *dockerRepository) PackageVersions(ch chan<- string, packageName string) (err error) {
	return docker.GetImageTags(ch, packageName, repo.connectionInfo.registry)
}

func (repo dockerRepository) Push(
	kpmHomeDir string,
	packageInfo *template_package.PackageInfo,
) error {
	return fmt.Errorf("not yet implemented")
}

func (repo *dockerRepository) Pull(
	kpmHomeDir string,
	packageInfo *template_package.PackageInfo,
) error {
	return fmt.Errorf("not yet implemented")
}

func repoInfoToDockerRepo(repoInfo *repositoryInfo) (Repository, error) {
	var err error

	var result = &dockerRepository{name: repoInfo.Name}

	var connectionInfo dockerRepositoryConnectionInfo
	err = repoInfo.ConnectionInfo.Decode(&connectionInfo)
	if err != nil {
		return result, fmt.Errorf("docker repository connection info is not a valid structure")
	}

	result.connectionInfo = connectionInfo

	return result, nil
}
