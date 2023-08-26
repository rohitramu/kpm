package template_repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/rohitramu/kpm/pkg/utils/templates"
)

// TODO: 1) Config for adding repos.  2) Commands for repo operations (list, pull, push).

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

func (repo *filesystemRepository) FindPackages(searchTerm string) (result []*templates.PackageInfo, err error) {
	var packageRepoDir = repo.absoluteFilePath

	var fullPackageNames []string
	fullPackageNames, err = template_package.GetPackageFullNamesFromLocalRepository(packageRepoDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to enumerate template packages in local repository named '%s': %s",
			repo.name,
			err,
		)
	}

	for _, fullPackageName := range fullPackageNames {
		var packageInfo *templates.PackageInfo
		packageInfo, err = template_package.GetPackageInfo(repo.absoluteFilePath, fullPackageName)
		if err != nil {
			return nil, err
		}

		// If the package name doesn't contain the search term, ignore it.
		if searchTerm != "" && !strings.Contains(packageInfo.Name, searchTerm) {
			continue
		}

		result = append(result, packageInfo)
	}

	return result, nil
}

func (repo *filesystemRepository) PackageVersions(packageName string) (result []string, err error) {
	return template_package.GetPackageVersions(repo.absoluteFilePath, packageName)
}

func (repo *filesystemRepository) Push(
	kpmHomeDir string,
	packageInfo *templates.PackageInfo,
) (err error) {
	if packageInfo == nil {
		log.Panicf("packageInfo is nil")
	}

	// Get the source directory.
	var packageDirSrc = template_package.GetPackageDir(
		kpmHomeDir,
		template_package.GetPackageFullName(packageInfo.Name, packageInfo.Version),
	)

	// Get the destination directory.
	var packageDirDst = template_package.GetPackageDir(
		repo.absoluteFilePath,
		template_package.GetPackageFullName(packageInfo.Name, packageInfo.Version),
	)

	return copyPackage(packageDirSrc, packageDirDst)
}

func (repo *filesystemRepository) Pull(
	kpmHomeDir string,
	packageInfo *templates.PackageInfo,
) (err error) {
	if packageInfo == nil {
		log.Panicf("packageInfo is nil")
	}

	// Get the source directory.
	var packageDirSrc = template_package.GetPackageDir(
		repo.absoluteFilePath,
		template_package.GetPackageFullName(packageInfo.Name, packageInfo.Version),
	)

	// If the directory doesn't exist, return an error.
	err = files.DirExists(packageDirSrc, packageInfo.Name)
	if err != nil {
		return errors.Join(PackageNotFoundError{PackageInfo: *packageInfo})
	}

	// Get the destination directory.
	var packageDirDst = template_package.GetPackageDir(
		kpmHomeDir,
		template_package.GetPackageFullName(packageInfo.Name, packageInfo.Version),
	)

	return copyPackage(packageDirSrc, packageDirDst)
}

func repoInfoToFilesystemRepo(repoInfo *RepositoryInfo) (Repository, error) {
	if repoInfo == nil {
		log.Panicf("repoInfo is nil")
	}

	var err error
	var result = &filesystemRepository{name: repoInfo.Name}

	dirPath, ok := repoInfo.Location.(string)
	if !ok {
		return result, fmt.Errorf("filesystem repository connection info is not an absolute directory path")
	}

	// Make sure it's an absolute path or rooted in the home directory.
	if !files.IsAbsFromHomeOrRoot(dirPath) {
		return result, fmt.Errorf("the provided filesystem repository path is not an absolute path")
	}

	var absDirPath string
	absDirPath, err = files.GetAbsolutePath(dirPath)
	if err != nil {
		return result, err
	}

	// Make sure the directory exists.
	if err = files.DirExists(absDirPath, fmt.Sprintf("%s repo", repoInfo.Name)); err != nil {
		return result, fmt.Errorf("the provided filesystem repository path does not point to a directory: %s", err)
	}

	// TODO: Make sure all folders inside this local repo are valid template packages.

	result.absoluteFilePath = absDirPath

	return result, nil
}

func copyPackage(packageDirSrc string, packageDirDst string) (err error) {
	// Delete the destination directory.
	if err = files.DeleteDirIfExists(packageDirDst, "destination template package", true); err != nil {
		return err
	}

	// Copy package to output directory.
	log.Debugf("Copying package from '%s' to '%s'", packageDirSrc, packageDirDst)
	err = files.CopyDir(packageDirSrc, packageDirDst)
	if err != nil {
		return err
	}

	return err
}