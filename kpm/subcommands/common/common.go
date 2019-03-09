package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/stacks/linkedliststack"

	"../utils/constants"
	"../utils/files"
	"../utils/log"
	"../utils/templates"
	"../utils/types"
	"../utils/validation"
	"../utils/yaml"
)

// PullPackage retrieves a remote template package and makes it available for use.  If a package
// was successfully retrieved, this function returns the retrieved version number.
func PullPackage(packageName string, wildcardPackageVersion string) (string, error) {
	//TODO: Get list of versions in remote repository

	//TODO: Resolve version to the highest that is compatible with the requested version

	//TODO: Download the template package of the resolved version into the local package repository

	return "", fmt.Errorf("Could not find version matching \"%s\" in remote repository for package: %s", wildcardPackageVersion, packageName)
}

// GetTemplateInput creates the input values for a template by combining the interface, parameters and package info.
func GetTemplateInput(parentTemplate *template.Template, packageDirPath string, parameters *types.GenericMap) (*types.GenericMap, error) {
	var err error

	// Add top-level objects
	var result = types.GenericMap{}
	result[constants.TemplateFieldPackage], err = GetPackageInfo(packageDirPath)
	if err != nil {
		return nil, err
	}
	result[constants.TemplateFieldValues], err = getValuesFromInterface(parentTemplate, packageDirPath, parameters)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSharedTemplate creates a template which contains default options, functions and
// helper template definitions defined in the given package.
func GetSharedTemplate(packageDirPath string) (*template.Template, error) {
	var err error

	// Get the directory which contains the helper templates
	var helpersDirPath = constants.GetHelpersDirPath(packageDirPath)

	// Create a template which includes the helper template definitions
	var sharedTemplate *template.Template
	var numHelpers int
	sharedTemplate, numHelpers, err = templates.ChainTemplatesFromDir(templates.GetRootTemplate(), helpersDirPath)
	if err != nil {
		return nil, err
	}

	log.Verbose(fmt.Sprintf("Found %d template(s) in directory: %s", numHelpers, helpersDirPath))

	return sharedTemplate, nil
}

// GetPackageInfo validates the package directory and returns the package info object for a given package.
func GetPackageInfo(packageDirPath string) (*types.PackageInfo, error) {
	var err error

	// Make sure that the package exists
	err = files.DirExists(packageDirPath, "package")
	if err != nil {
		return nil, err
	}

	// Check that the package info file exists
	var packageInfoFilePath = filepath.Join(packageDirPath, constants.PackageInfoFileName)
	err = files.FileExists(packageInfoFilePath, "package information")
	if err != nil {
		return nil, err
	}

	// Get package info file content
	var yamlBytes []byte
	yamlBytes, err = files.ReadBytes(packageInfoFilePath)
	if err != nil {
		return nil, err
	}

	// Get package info object from file content
	var packageInfo = new(types.PackageInfo)
	err = yaml.BytesToObject(yamlBytes, packageInfo)
	if err != nil {
		return nil, err
	}

	// Validate package name
	err = validation.ValidatePackageName(packageInfo.Name)
	if err != nil {
		return nil, err
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageInfo.Version, false)
	if err != nil {
		return nil, err
	}

	// Make sure that the interface file exists
	var interfaceFilePath = filepath.Join(packageDirPath, constants.InterfaceFileName)
	err = files.FileExists(interfaceFilePath, "interface")
	if err != nil {
		return nil, err
	}

	// Make sure that the parameters file exists
	var parametersFilePath = constants.GetDefaultParametersFilePath(packageDirPath)
	err = files.FileExists(parametersFilePath, "parameters")
	if err != nil {
		return nil, err
	}

	return packageInfo, nil
}

// GetPackageParameters returns the parameters in a file as an object which can be used as input to the interface template in a package.
func GetPackageParameters(parametersFilePath string) (*types.GenericMap, error) {
	var err error

	// Make sure that the parameters file exists
	err = files.FileExists(parametersFilePath, "parameters")
	if err != nil {
		return nil, err
	}

	// Get parameters file content as bytes
	var parametersFileBytes []byte
	parametersFileBytes, err = files.ReadBytes(parametersFilePath)
	if err != nil {
		return nil, err
	}

	// Get parameters
	var parameters = new(types.GenericMap)
	err = yaml.BytesToObject(parametersFileBytes, parameters)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}

// GetExecutableTemplates returns all executable templates in a template package.
func GetExecutableTemplates(parentTemplate *template.Template, packageDirPath string) ([]*template.Template, error) {
	var err error

	// Get the templates directory
	var executableTemplatesDir = constants.GetTemplatesDirPath(packageDirPath)
	err = files.DirExists(executableTemplatesDir, "templates")
	if err != nil {
		return nil, err
	}

	// Return the templates in the directory
	log.Verbose(fmt.Sprintf("Found template directory: %s", executableTemplatesDir))
	var result []*template.Template
	result, err = templates.GetTemplatesFromDir(parentTemplate, executableTemplatesDir)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetDependencyDefinitionTemplates returns the templates for all dependency definition templates in a template package.
func GetDependencyDefinitionTemplates(parentTemplate *template.Template, packageDirPath string) ([]*template.Template, error) {
	var err error

	// Get the dependencies directory
	var dependenciesDir = constants.GetDependenciesDirPath(packageDirPath)
	err = files.DirExists(dependenciesDir, "dependencies")
	if err != nil {
		return nil, err
	}

	var dependencyTemplates []*template.Template
	dependencyTemplates, err = templates.GetTemplatesFromDir(parentTemplate, dependenciesDir)
	if err != nil {
		return nil, err
	}

	return dependencyTemplates, nil
}

// getValuesFromInterface creates the values which can be used as input to templates by executing the interface with parameters.
func getValuesFromInterface(parentTemplate *template.Template, packageDirPath string, parameters *types.GenericMap) (*types.GenericMap, error) {
	var err error

	// Create template object from interface file
	var templateName = constants.InterfaceFileName
	var interfaceFilePath = filepath.Join(packageDirPath, templateName)
	var tmpl *template.Template
	tmpl, err = templates.GetTemplateFromFile(parentTemplate, templateName, interfaceFilePath)
	if err != nil {
		return nil, err
	}

	// Generate values by applying parameters to interface
	var interfaceBytes []byte
	interfaceBytes, err = templates.ExecuteTemplate(tmpl, parameters)
	if err != nil {
		return nil, err
	}

	// Get values object from generated values yaml file
	var result = new(types.GenericMap)
	err = yaml.BytesToObject(interfaceBytes, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetPackageNamesFromLocalRepository returns the list of package names in the local KPM package repository.
func GetPackageNamesFromLocalRepository(packageRepositoryDir string) ([]string, error) {
	var err error
	var ok bool

	// Exit early if the packages directory doesn't exist
	err = files.DirExists(packageRepositoryDir, "packages repository")
	if err != nil {
		return nil, err
	}

	// Traverse the packages directory
	var packages = treeset.NewWithStringComparator()
	var toVisit = linkedliststack.New()
	toVisit.Push(packageRepositoryDir)
	for currentPathObj, ok := toVisit.Pop(); ok; currentPathObj, ok = toVisit.Pop() {
		// Assert type as string
		var currentPath string
		currentPath, ok = currentPathObj.(string)
		if !ok {
			// We should never fail here since we are providing the values
			log.Panic(fmt.Sprintf("Unexpected object when string was expected: %s", reflect.TypeOf(currentPathObj)))
		}

		// Get the file info
		var fileInfo os.FileInfo
		fileInfo, err = os.Stat(currentPath)
		if err != nil {
			// We should never fail here since we are providing the values
			log.Panic(err)
		}

		// Ignore files
		if !fileInfo.IsDir() {
			continue
		}

		// Check if this is a valid package directory
		_, err = GetPackageInfo(currentPath)
		if err == nil {
			// Found a valid package, so add it to the list of found packages
			packages.Add(currentPath)

			// We don't want to add subdirectories, since packages cannot be nested
			continue
		}

		// Since this is not a valid package directory, visit all subdirectories if this is the root directory or the directory's name is a valid namespace name
		if currentPath == packageRepositoryDir || validation.ValidateNamespaceSegment(filepath.Base(currentPath)) == nil {
			var subdirectories []os.FileInfo
			subdirectories, err = ioutil.ReadDir(currentPath)
			if err != nil {
				return nil, err
			}
			for _, dir := range subdirectories {
				var dirName = dir.Name()
				if dir.IsDir() {
					toVisit.Push(filepath.Join(currentPath, dirName))
				}
			}
		}
	}

	// Compile the list of results
	var results = make([]string, packages.Size())
	for it := packages.Iterator(); it.Next(); {
		var path string
		path, ok = it.Value().(string)
		if !ok {
			log.Panic(fmt.Sprintf("Unexpected type found when getting list of string package names: %s", reflect.TypeOf(it.Value())))
		}

		// Get the relative path
		path, err = filepath.Rel(packageRepositoryDir, path)
		if err != nil {
			return nil, err
		}

		// Get the image name by using forward slashes instead of backward slashes (if on a Windows machine)
		var imageName = filepath.ToSlash(path)
		results[it.Index()] = imageName
	}

	return results, nil
}
