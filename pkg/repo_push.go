package pkg

import (
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/rohitramu/kpm/pkg/utils/template_repository"
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

func PushToRepository(
	kpmHomeDir string,
	repos *template_repository.RepositoryCollection,
	repoName string,
	packageName string,
	packageVersion string,
	userHasConfirmed bool,
) (err error) {
	if packageVersion == "" {
		packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName)
		if err != nil {
			return err
		}
	}

	var packageInfo = &templates.PackageInfo{
		Name:    packageName,
		Version: packageVersion,
	}

	return repos.Push(kpmHomeDir, repoName, packageInfo, userHasConfirmed)
}
