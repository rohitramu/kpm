package template_repository

import (
	"fmt"
	"path/filepath"

	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

const repositoryTypeNameFilesystem = "local"

type filesystemRepository struct {
	name             string
	absoluteFilePath string
}

func (repo *filesystemRepository) GetName() string {
	return repo.name
}

func (repo *filesystemRepository) GetType() string {
	return repositoryTypeNameFilesystem
}

func (repo *filesystemRepository) Packages() ([]templates.PackageInfo, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (repo *filesystemRepository) PackageVersions() ([]string, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (repo filesystemRepository) Push(templates.PackageInfo) error {
	return fmt.Errorf("not yet implemented")
}

func (repo *filesystemRepository) Pull(templates.PackageInfo) error {
	return fmt.Errorf("not yet implemented")
}

func repoInfoToFilesystemRepo(repoInfo RepositoryInfo) (Repository, error) {
	var result = &filesystemRepository{name: repoInfo.Name}

	dirPath, ok := repoInfo.ConnectionInfo.(string)
	if !ok {
		return result, fmt.Errorf("filesystem repository connection info is not a string directory path")
	}

	// Make sure it's an absolute path.
	if !filepath.IsAbs(dirPath) {
		return result, fmt.Errorf("the provided filesystem repository path is not an absolute path")
	}

	// Make sure the directory exists.
	if err := files.DirExists(dirPath, fmt.Sprintf("%s repo", repoInfo.Name)); err != nil {
		return result, fmt.Errorf("the provided filesystem repository path does not point to a directory: %s", err)
	}

	// TODO: Make sure all folders inside this local repo are valid template packages.

	result.absoluteFilePath = dirPath

	return result, nil
}
