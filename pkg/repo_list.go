package pkg

import "github.com/rohitramu/kpm/pkg/utils/template_repository"

func ListPackageRepositories(
	repos *template_repository.RepositoryCollection,
) ([]string, error) {
	return repos.GetRepositoryNames(), nil
}
