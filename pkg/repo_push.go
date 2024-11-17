package pkg

import (
	"errors"

	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/rohitramu/kpm/pkg/utils/template_repository"
	"github.com/rohitramu/kpm/pkg/utils/user_prompts"
)

func PushPackage(
	kpmHomeDir string,
	repos *template_repository.RepositoryCollection,
	repoName string,
	packageName string,
	packageVersion string,
	userHasConfirmed bool,
) (err error) {
	var repoNames = repos.GetRepositoryNames()
	if len(repoNames) == 0 {
		return errors.New("no repositories configured")
	}

	if packageVersion == "" {
		packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName)
		if err != nil {
			return err
		}
	}

	var packageInfo = &template_package.PackageInfo{
		Name:    packageName,
		Version: packageVersion,
	}

	if repoName == "" {
		repoName = repoNames[0]
	}

	if !userHasConfirmed {
		user_prompts.ConfirmWithUser(
			"Pushing package '%s' (version '%s') to primary repository '%s'",
			packageName,
			packageVersion,
			repoName,
		)
	}

	var repo template_repository.Repository
	repo, err = repos.GetRepository(repoName)
	if err != nil {
		return err
	}

	return repo.Push(kpmHomeDir, packageInfo)
}
