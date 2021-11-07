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

	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/docker"
	"github.com/rohitramu/kpm/subcommands/utils/files"
	"github.com/rohitramu/kpm/subcommands/utils/log"
	"github.com/rohitramu/kpm/subcommands/utils/templates"
	"github.com/rohitramu/kpm/subcommands/utils/types"
	"github.com/rohitramu/kpm/subcommands/utils/validation"
	"github.com/rohitramu/kpm/subcommands/utils/yaml"
)

// PullPackage retrieves a remote template package and makes it available for use.
func PullPackage(kpmHomeDir string, dockerRegistry string, packageName string, packageVersion string) error {
	//TODO: Decouple Docker operations from the rest of KPM.
	var err error

	// Get the package's full name
	var packageFullName = constants.GetPackageFullName(packageName, packageVersion)

	log.Info("Pulling from \"%s\": %s", dockerRegistry, packageFullName)

	// Get the image name
	var imageName = docker.GetImageName(dockerRegistry, packageName, packageVersion)

	// Pull the Docker image
	err = docker.PullImage(imageName)
	if err != nil {
		return err
	}

	// Delete the local image after we're done
	defer func() {
		var deleteErr = docker.DeleteImage(imageName)
		if deleteErr != nil {
			if err != nil {
				err = fmt.Errorf("Failed to delete image: %s\n%s\n%s", imageName, deleteErr, err)
			}
		}
	}()

	// Get the package directory
	var packageDir = constants.GetPackageDir(kpmHomeDir, packageFullName)

	// Extract Docker image contents into the local package repository
	err = docker.ExtractImageContents(imageName, packageDir)
	if err != nil {
		log.Error("Failed to extract Docker image contents")
		return err
	}

	//TODO: Pull all dependencies of this package

	return nil
}

// GetTemplateInput creates the input values for a template by combining the interface, parameters and package info.
func GetTemplateInput(kpmHomeDir string, packageFullName string, parentTemplate *template.Template, parameters *types.GenericMap) (*types.GenericMap, error) {
	var err error

	var packageDir = constants.GetPackageDir(kpmHomeDir, packageFullName)
	var result = types.GenericMap{}

	// Add package info
	var packageInfo *types.PackageInfo
	packageInfo, err = GetPackageInfo(kpmHomeDir, packageDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to get information about package: %s\n%s", packageFullName, err)
	}
	var packageInfoMap = types.GenericMap{}
	packageInfoMap["name"] = packageInfo.Name
	packageInfoMap["version"] = packageInfo.Version
	result[constants.TemplateFieldPackage] = &packageInfoMap

	// Get the default values
	var inputParameters *types.GenericMap
	inputParameters, err = GetPackageParameters(constants.GetDefaultParametersFile(packageDir))
	if err != nil {
		return nil, err
	}

	// If the file didn't exist, create an empty map
	if inputParameters == nil {
		inputParameters = new(types.GenericMap)
	}

	// Allow default values to be overridden by the provided parameters
	for key := range *parameters {
		(*inputParameters)[key] = (*parameters)[key]
	}

	// Add values
	result[constants.TemplateFieldValues], err = getValuesFromInterface(parentTemplate, packageDir, inputParameters)
	if err != nil {
		return nil, fmt.Errorf("Failed to calculate values from the interface in package: %s\n%s", packageFullName, err)
	}

	return &result, nil
}

// GetSharedTemplate creates a template which contains default options, functions and
// helper template definitions defined in the given package.
func GetSharedTemplate(packageDir string) (*template.Template, error) {
	var err error

	// Get the directory which contains the helper templates
	var helpersDir = constants.GetHelpersDir(packageDir)

	// Get the root template
	var sharedTemplate = templates.NewRootTemplate()

	// Create a template which includes the helper template definitions
	if files.DirExists(helpersDir, "helpers") == nil {
		var numHelpers int
		sharedTemplate, numHelpers, err = templates.ChainTemplatesFromDir(sharedTemplate, helpersDir)
		if err != nil {
			return nil, err
		}

		log.Verbose("Found %d template(s) in directory: %s", numHelpers, helpersDir)
	}

	// Add the package-specific template functions
	sharedTemplate = templates.AddPackageSpecificTemplateFunctions(sharedTemplate)

	return sharedTemplate, nil
}

// GetPackageInfo validates the package directory and returns the package info object for a given package.
func GetPackageInfo(kpmHomeDir string, packageDir string) (*types.PackageInfo, error) {
	var err error

	// Make sure that the package exists
	err = files.DirExists(packageDir, "package")
	if err != nil {
		return nil, err
	}

	// Check that the package info file exists
	var packageInfoFile = constants.GetPackageInfoFile(packageDir)
	err = files.FileExists(packageInfoFile, "package information")
	if err != nil {
		return nil, err
	}

	// Get package info file content
	var yamlBytes []byte
	yamlBytes, err = files.ReadBytes(packageInfoFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read package information file: %s\n%s", packageInfoFile, err)
	}

	// Get package info object from file content
	var packageInfo = new(types.PackageInfo)
	err = yaml.BytesToObject(yamlBytes, packageInfo)
	if err != nil {
		return nil, fmt.Errorf("Invalid package information file: %s\n%s", packageInfoFile, err)
	}

	// Validate package name
	var packageName = packageInfo.Name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return nil, err
	}

	// Validate package version
	var packageVersion = packageInfo.Version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return nil, err
	}

	// Make sure that the interface file exists
	var interfaceFilePath = constants.GetInterfaceFile(packageDir)
	err = files.FileExists(interfaceFilePath, "interface")
	if err != nil {
		return nil, err
	}

	// Make sure that the parameters file exists
	var parametersFile = constants.GetDefaultParametersFile(packageDir)
	err = files.FileExists(parametersFile, "default parameters")
	if err != nil {
		return nil, err
	}

	// Validate the templates directory if it exists
	var templatesDir = constants.GetTemplatesDir(packageDir)
	if files.DirExists(templatesDir, "templates") == nil {
		var fileInfos []os.FileInfo
		fileInfos, err = ioutil.ReadDir(templatesDir)
		if err != nil {
			// We already checked that this directory exists, so we shouldn't ever get to here
			log.Panic("Failed to read directory: %s\n%s", templatesDir, err)
		}

		for _, fileInfo := range fileInfos {
			// Get file name
			var fileName = fileInfo.Name()

			// Don't allow directories
			if fileInfo.IsDir() {
				return nil, fmt.Errorf("Directories are not allowed in the \"%s\" directory: %s", constants.TemplatesDirName, fileName)
			}
		}
	}

	// Validate the helpers directory if it exists
	var helpersDir = constants.GetHelpersDir(packageDir)
	if files.DirExists(helpersDir, "helpers") == nil {
		// Make sure all helper template files have the extension ".tpl"
		var fileInfos []os.FileInfo
		fileInfos, err = ioutil.ReadDir(helpersDir)
		if err != nil {
			// We already checked that this directory exists, so we shouldn't ever get to here
			log.Panic("Failed to read directory: %s\n%s", helpersDir, err)
		}

		for _, fileInfo := range fileInfos {
			// Get file name
			var fileName = fileInfo.Name()

			// Don't allow directories
			if fileInfo.IsDir() {
				return nil, fmt.Errorf("Directories are not allowed in the \"%s\" directory: %s", constants.HelpersDirName, fileName)
			}

			// Check file extension
			var validExtension = ".tpl"
			if filepath.Ext(fileName) != validExtension {
				return nil, fmt.Errorf("Invalid helpers - helpers files must be valid template files with the extension \"%s\": %s", validExtension, fileName)
			}
		}
	}

	// Validate the dependencies directory if it exists
	var dependenciesDir = constants.GetDependenciesDir(packageDir)
	if files.DirExists(dependenciesDir, "dependencies") == nil {
		// Make sure all dependencies files have the extension ".yaml"
		var fileInfos []os.FileInfo
		fileInfos, err = ioutil.ReadDir(dependenciesDir)
		if err != nil {
			// We already checked that this directory exists, so we shouldn't ever get to here
			log.Panic("Failed to read directory: %s\n%s", dependenciesDir, err)
		}

		for _, fileInfo := range fileInfos {
			// Get file name
			var fileName = fileInfo.Name()

			// Don't allow directories
			if fileInfo.IsDir() {
				return nil, fmt.Errorf("Directories are not allowed in the \"%s\" directory: %s", constants.DependenciesDirName, fileName)
			}

			// Check file extension
			var validExtension = ".yaml"
			if filepath.Ext(fileName) != validExtension {
				return nil, fmt.Errorf("Invalid dependency definition - dependency definition files must be valid yaml files with the extension \"%s\": %s", validExtension, fileName)
			}
		}
	}

	return packageInfo, nil
}

// GetPackageParameters returns the parameters in a file as an object which can be used as input to the interface template in a package.
func GetPackageParameters(parametersFile string) (*types.GenericMap, error) {
	var err error

	// Make sure that the parameters file exists
	err = files.FileExists(parametersFile, "parameters")
	if err != nil {
		return nil, err
	}

	// Get parameters file content as bytes
	var parametersFileBytes []byte
	parametersFileBytes, err = files.ReadBytes(parametersFile)
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
func GetExecutableTemplates(parentTemplate *template.Template, packageDir string) ([]*template.Template, error) {
	var err error

	// Get the templates directory
	var executableTemplatesDir = constants.GetTemplatesDir(packageDir)

	// If the templates directory doesn't exist, just return a list of no templates instead of erroring out
	if files.DirExists(executableTemplatesDir, "templates") != nil {
		return []*template.Template{}, nil
	}

	// Return the templates in the directory
	log.Verbose("Found template directory: %s", executableTemplatesDir)
	var result []*template.Template
	result, err = templates.GetTemplatesFromDir(parentTemplate, executableTemplatesDir)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetDependencyDefinitionTemplates returns the templates for all dependency definition templates in a template package.
func GetDependencyDefinitionTemplates(parentTemplate *template.Template, packageDir string) ([]*template.Template, error) {
	var err error

	// Get the dependencies directory
	var dependenciesDir = constants.GetDependenciesDir(packageDir)

	// If the dependencies directory doesn't exist, just return a list of no templates instead of erroring out
	if files.DirExists(dependenciesDir, "dependencies") != nil {
		return []*template.Template{}, nil
	}

	var dependencyTemplates []*template.Template
	dependencyTemplates, err = templates.GetTemplatesFromDir(parentTemplate, dependenciesDir)
	if err != nil {
		return nil, err
	}

	return dependencyTemplates, nil
}

// getValuesFromInterface creates the values which can be used as input to templates by executing the interface with parameters.
func getValuesFromInterface(parentTemplate *template.Template, packageDir string, parameters *types.GenericMap) (*types.GenericMap, error) {
	var err error

	// Create template object from interface file
	var interfaceFile = constants.GetInterfaceFile(packageDir)
	var tmpl *template.Template
	tmpl, err = templates.GetTemplateFromFile(parentTemplate, filepath.Base(interfaceFile), interfaceFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to get interface file: %s\n%s", packageDir, err)
	}

	// Generate values by applying parameters to interface
	var interfaceBytes []byte
	interfaceBytes, err = templates.ExecuteTemplate(tmpl, parameters)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute interface file: %s\n%s", interfaceFile, err)
	}

	// Get values object from generated values yaml file
	var result = new(types.GenericMap)
	err = yaml.BytesToObject(interfaceBytes, result)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse generated values from interface file: %s\n%s", interfaceFile, err)
	}

	return result, nil
}

// GetPackageFullNamesFromLocalRepository returns the list of package names in the local KPM package repository, in alphabetical order.
func GetPackageFullNamesFromLocalRepository(kpmHomeDir string) ([]string, error) {
	var err error
	var ok bool

	var packageRepositoryDir = constants.GetPackageRepositoryDir(kpmHomeDir)

	// Exit early if the packages directory doesn't exist
	err = files.DirExists(packageRepositoryDir, "packages repository")
	if err != nil {
		return []string{}, nil
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
			log.Panic("Unexpected object when string was expected: %s", reflect.TypeOf(currentPathObj))
		}

		// Get the file info
		var fileInfo os.FileInfo
		fileInfo, err = os.Stat(currentPath)
		if err != nil {
			// We should never fail here since we are providing the values
			log.Panic("Unexpected file path: %s", err)
		}

		// Ignore files
		if !fileInfo.IsDir() {
			continue
		}

		// Check if this is a valid package directory
		var packageInfo *types.PackageInfo
		packageInfo, err = GetPackageInfo(kpmHomeDir, currentPath)
		if err == nil {
			// Get the package's full name
			var packageFullName = constants.GetPackageFullName(packageInfo.Name, packageInfo.Version)

			// Calculate what the full name of the package should be
			var packageFullNameFromPath string
			packageFullNameFromPath, err = filepath.Rel(packageRepositoryDir, currentPath)
			if err != nil {
				log.Panic("Failed to get relative path: %s -> %s", packageRepositoryDir, currentPath)
			}

			// We always expect forward slashes for namespaces
			packageFullNameFromPath = filepath.ToSlash(packageFullNameFromPath)

			// Check that the name of the directory matches the package's full name
			if packageFullNameFromPath != packageFullName {
				// Log a warning
				log.Warning("Found corrupted package in local repository (directory name does not match package name): %s", currentPath)

				// Don't return this package, just continue looking for other packages
				continue
			}

			// Found a valid package, so add it to the list of found packages
			packages.Add(packageFullNameFromPath)

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

	var result = make([]string, packages.Size())
	var i = 0
	for it := packages.Iterator(); it.Next(); {
		var packageName string
		var value = it.Value()
		packageName, ok = value.(string)
		if !ok {
			log.Panic("Unexpected value in tree map: %s", value)
		}

		result[i] = packageName
		i++
	}

	return result, nil
}
