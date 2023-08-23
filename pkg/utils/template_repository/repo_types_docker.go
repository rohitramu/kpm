package template_repository

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

const repositoryTypeNameDocker = "docker"

type dockerRepository struct {
	name           string
	connectionInfo dockerRepositoryConnectionInfo
}

type dockerRepositoryConnectionInfo struct {
	dockerUsername     string
	dockerOrganization string
}

func (repo *dockerRepository) GetName() string {
	return repo.name
}

func (repo *dockerRepository) GetType() string {
	return repositoryTypeNameFilesystem
}

func (repo *dockerRepository) Packages() ([]*templates.PackageInfo, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (repo *dockerRepository) PackageVersions(packageName string) ([]string, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (repo dockerRepository) Push(
	kpmHomeDir string,
	packageInfo *templates.PackageInfo,
	userHasConfirmed bool,
) error {
	return fmt.Errorf("not yet implemented")
}

func (repo *dockerRepository) Pull(
	kpmHomeDir string,
	packageInfo *templates.PackageInfo,
	userHasConfirmed bool,
) error {
	return fmt.Errorf("not yet implemented")
}

func repoInfoToDockerRepo(repoInfo *RepositoryInfo) (Repository, error) {
	var result = &dockerRepository{name: repoInfo.Name}

	connectionInfo, ok := repoInfo.Location.(dockerRepositoryConnectionInfo)
	if !ok {
		return result, fmt.Errorf("docker repository connection info is not a valid structure")
	}

	log.Warningf("Docker connection info: %v", connectionInfo)

	return result, nil
}
