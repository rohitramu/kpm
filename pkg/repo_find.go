package pkg

import (
	"errors"

	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/rohitramu/kpm/pkg/utils/template_repository"
)

func FindPackages(
	ch chan<- *template_package.PackageInfo,
	kpmHomeDir string,
	repos *template_repository.RepositoryCollection,
	repoName string,
	searchTerm string,
) (err error) {
	var repoNames = repos.GetRepositoryNames()
	if len(repoNames) == 0 {
		return errors.New("no repositories configured")
	}

	for _, repoName := range repos.GetRepositoryNames() {
		var repo template_repository.Repository
		repo, err = repos.GetRepository(repoName)
		if err != nil {
			return err
		}

		err = repo.FindPackages(ch, searchTerm)
		if err != nil && !errors.Is(err, template_repository.PackageNotFoundError{}) {
			return err
		}
	}

	return nil
}
