package pkg

import (
	"errors"

	"github.com/rohitramu/kpm/src/pkg/utils/template_package"
	"github.com/rohitramu/kpm/src/pkg/utils/template_repository"
)

func PullPackage(
	kpmHomeDir string,
	repos *template_repository.RepositoryCollection,
	repoName string,
	packageName string,
	packageVersion string,
) (err error) {
	var repoNames = repos.GetRepositoryNames()
	if len(repoNames) == 0 {
		return errors.New("no repositories configured")
	}

	var packageInfo = &template_package.PackageInfo{
		Name:    packageName,
		Version: packageVersion,
	}

	// Repo name was provided, so only pull from the repo.
	if repoName != "" {
		var repo template_repository.Repository
		repo, err = repos.GetRepository(repoName)
		if err != nil {
			return err
		}

		return repo.Pull(kpmHomeDir, packageInfo)
	}

	// If repo name wasn't provided, check all repos for the package.
	for _, repoName := range repoNames {
		var repo template_repository.Repository
		repo, err = repos.GetRepository(repoName)
		if err != nil {
			return err
		}

		err = repo.Pull(kpmHomeDir, packageInfo)
		if err == nil {
			// Found the package, so return.
			return nil
		}

		// If the package wasn't found, continue looking in other repos.
		if errors.Is(err, template_repository.PackageNotFoundError{}) {
			continue
		}

		return err
	}

	return template_repository.PackageNotFoundError{PackageInfo: *packageInfo}
}
