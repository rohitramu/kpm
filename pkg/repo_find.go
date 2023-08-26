package pkg

import (
	"errors"

	"github.com/rohitramu/kpm/pkg/utils/template_repository"
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

func FindPackages(
	kpmHomeDir string,
	repos *template_repository.RepositoryCollection,
	repoName string,
	searchTerm string,
) (result []*templates.PackageInfo, err error) {
	result = make([]*templates.PackageInfo, 0)

	var repoNames = repos.GetRepositoryNames()
	if len(repoNames) == 0 {
		return nil, errors.New("no repositories configured")
	}

	for _, repoName := range repos.GetRepositoryNames() {
		var repo template_repository.Repository
		repo, err = repos.GetRepository(repoName)
		if err != nil {
			return nil, err
		}

		var packageInfos []*templates.PackageInfo
		packageInfos, err = repo.FindPackages(searchTerm)
		if err != nil && !errors.Is(err, template_repository.PackageNotFoundError{}) {
			return nil, err
		}

		result = append(result, packageInfos...)
	}

	return result, nil
}
