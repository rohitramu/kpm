package template_package

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/emirpasic/gods/sets/treeset"

	"github.com/rohitramu/kpm/src/pkg/utils/constants"
	"github.com/rohitramu/kpm/src/pkg/utils/files"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
	"github.com/rohitramu/kpm/src/pkg/utils/templates"
	"github.com/rohitramu/kpm/src/pkg/utils/validation"
	"github.com/rohitramu/kpm/src/pkg/utils/yaml"
)

// TODO: Move repo-specific methods to the "template_repository" package.

// GetDefaultParametersFile returns the path of the default parameters file in a template package.
func GetDefaultParametersFile(packageDir string) string {
	var parametersFilePath = filepath.Join(packageDir, constants.ParametersFileName)

	return parametersFilePath
}

// GetInterfaceFile returns the path of the interface file in a template package.
func GetInterfaceFile(packageDir string) string {
	var interfaceFilePath = filepath.Join(packageDir, constants.InterfaceFileName)

	return interfaceFilePath
}

// GetPackageInfoFile returns the path of the package information file in a template package.
func GetPackageInfoFile(packageDir string) string {
	var packageInfoFilePath = filepath.Join(packageDir, constants.PackageInfoFileName)

	return packageInfoFilePath
}

// GetPackageFullName returns the full package name with version.
func GetPackageFullName(packageName string, packageVersion string) string {
	return fmt.Sprintf("%s-%s", packageName, packageVersion)
}

// GetTemplateInput creates the input values for a template by combining the interface, parameters and package info.
func GetTemplateInput(
	kpmHomeDir string,
	packageFullName string,
	parentTemplate *template.Template,
	parameters *map[string]any,
) (*map[string]any, error) {
	var err error

	var packageDir = GetPackageDir(kpmHomeDir, packageFullName)
	var result = map[string]any{}

	// Add package info
	var packageInfo *PackageInfo
	packageInfo, err = GetPackageInfo(packageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get information about package: %s\n%s", packageFullName, err)
	}
	var packageInfoMap = map[string]any{}
	packageInfoMap["name"] = packageInfo.Name
	packageInfoMap["version"] = packageInfo.Version
	result[constants.TemplateFieldPackage] = &packageInfoMap

	// Get the default values
	var inputParameters *map[string]any
	inputParameters, err = GetPackageParameters(GetDefaultParametersFile(packageDir))
	if err != nil {
		return nil, err
	}

	// If the file didn't exist, create an empty map
	if inputParameters == nil {
		inputParameters = new(map[string]any)
	}

	// Allow default values to be overridden by the provided parameters
	for key := range *parameters {
		(*inputParameters)[key] = (*parameters)[key]
	}

	// Add values
	result[constants.TemplateFieldValues], err = getValuesFromInterface(parentTemplate, packageDir, inputParameters)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate values from the interface in package: %s\n%s", packageFullName, err)
	}

	return &result, nil
}

// GetSharedTemplate creates a template which contains default options, functions and
// helper template definitions defined in the given package.
func GetSharedTemplate(packageDir string) (*template.Template, error) {
	var err error

	// Get the directory which contains the helper templates
	var helpersDir = GetHelpersDir(packageDir)

	// Get the root template
	var sharedTemplate = templates.NewRootTemplate()

	// Create a template which includes the helper template definitions
	if files.DirExists(helpersDir, "helpers") == nil {
		var numHelpers int
		sharedTemplate, numHelpers, err = templates.ChainTemplatesFromDir(sharedTemplate, helpersDir)
		if err != nil {
			return nil, err
		}

		log.Debugf("Found %d template(s) in directory: %s", numHelpers, helpersDir)
	}

	// Add the package-specific template functions
	sharedTemplate = templates.AddPackageSpecificTemplateFunctions(sharedTemplate)

	return sharedTemplate, nil
}

// GetPackageInfo validates the package directory and returns the package info object for a given package.
func GetPackageInfo(packageDirAbsPath string) (*PackageInfo, error) {
	var err error

	// Make sure that the package exists
	err = files.DirExists(packageDirAbsPath, "package")
	if err != nil {
		return nil, err
	}

	// Check that the package info file exists
	var packageInfoFile = GetPackageInfoFile(packageDirAbsPath)
	err = files.FileExists(packageInfoFile, "template package information")
	if err != nil {
		return nil, err
	}

	// Get package info file content
	var yamlBytes []byte
	yamlBytes, err = files.ReadBytes(packageInfoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read template package information file: %s\n%s", packageInfoFile, err)
	}

	// Get package info object from file content
	var packageInfo = new(PackageInfo)
	err = yaml.BytesToObject(yamlBytes, packageInfo)
	if err != nil {
		return nil, fmt.Errorf("invalid package information file: %s\n%s", packageInfoFile, err)
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
	var interfaceFilePath = GetInterfaceFile(packageDirAbsPath)
	err = files.FileExists(interfaceFilePath, "interface")
	if err != nil {
		return nil, err
	}

	// Make sure that the parameters file exists
	var parametersFile = GetDefaultParametersFile(packageDirAbsPath)
	err = files.FileExists(parametersFile, "default parameters")
	if err != nil {
		return nil, err
	}

	// Validate the templates directory if it exists
	var templatesDir = GetTemplatesDir(packageDirAbsPath)
	if files.DirExists(templatesDir, "templates") == nil {
		var fileInfos []os.DirEntry
		fileInfos, err = os.ReadDir(templatesDir)
		if err != nil {
			// We already checked that this directory exists, so we shouldn't ever get to here
			log.Panicf("Failed to read directory: %s\n%s", templatesDir, err)
		}

		for _, fileInfo := range fileInfos {
			// Get file name
			var fileName = fileInfo.Name()

			// Don't allow directories
			if fileInfo.IsDir() {
				return nil, fmt.Errorf("directories are not allowed in the \"%s\" directory: %s", constants.TemplatesDirName, fileName)
			}
		}
	}

	// Validate the helpers directory if it exists
	var helpersDir = GetHelpersDir(packageDirAbsPath)
	if files.DirExists(helpersDir, "helpers") == nil {
		// Make sure all helper template files have the extension ".tpl"
		var fileInfos []os.DirEntry
		fileInfos, err = os.ReadDir(helpersDir)
		if err != nil {
			// We already checked that this directory exists, so we shouldn't ever get to here
			log.Panicf("Failed to read directory: %s\n%s", helpersDir, err)
		}

		for _, fileInfo := range fileInfos {
			// Get file name
			var fileName = fileInfo.Name()

			// Don't allow directories
			if fileInfo.IsDir() {
				return nil, fmt.Errorf("directories are not allowed in the \"%s\" directory: %s", constants.HelpersDirName, fileName)
			}

			// Check file extension
			var validExtension = ".tpl"
			if filepath.Ext(fileName) != validExtension {
				return nil, fmt.Errorf("invalid helpers - helpers files must be valid template files with the extension \"%s\": %s", validExtension, fileName)
			}
		}
	}

	// Validate the dependencies directory if it exists
	var dependenciesDir = GetDependenciesDir(packageDirAbsPath)
	if files.DirExists(dependenciesDir, "dependencies") == nil {
		// Make sure all dependencies files have the extension ".yaml"
		var fileInfos []os.DirEntry
		fileInfos, err = os.ReadDir(dependenciesDir)
		if err != nil {
			// We already checked that this directory exists, so we shouldn't ever get to here
			log.Panicf("Failed to read directory: %s\n%s", dependenciesDir, err)
		}

		for _, fileInfo := range fileInfos {
			// Get file name
			var fileName = fileInfo.Name()

			// Don't allow directories
			if fileInfo.IsDir() {
				return nil, fmt.Errorf("directories are not allowed in the \"%s\" directory: %s", constants.DependenciesDirName, fileName)
			}

			// Check file extension
			var validExtension = ".yaml"
			if filepath.Ext(fileName) != validExtension {
				return nil, fmt.Errorf("invalid dependency definition - dependency definition files must be valid yaml files with the extension \"%s\": %s", validExtension, fileName)
			}
		}
	}

	return packageInfo, nil
}

// GetPackageParameters returns the parameters in a file as an object which can be used as input to the interface template in a package.
func GetPackageParameters(parametersFile string) (*map[string]any, error) {
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
	var parameters = new(map[string]any)
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
	var executableTemplatesDir = GetTemplatesDir(packageDir)

	// If the templates directory doesn't exist, just return a list of no templates instead of erroring out
	if files.DirExists(executableTemplatesDir, "templates") != nil {
		return []*template.Template{}, nil
	}

	// Return the templates in the directory
	log.Debugf("Found template directory: %s", executableTemplatesDir)
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
	var dependenciesDir = GetDependenciesDir(packageDir)

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
func getValuesFromInterface(
	parentTemplate *template.Template,
	packageDir string,
	parameters *map[string]any,
) (*map[string]any, error) {
	var err error

	// Create template object from interface file
	var interfaceFile = GetInterfaceFile(packageDir)
	var tmpl *template.Template
	tmpl, err = templates.GetTemplateFromFile(parentTemplate, filepath.Base(interfaceFile), interfaceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get interface file: %s\n%s", packageDir, err)
	}

	// Generate values by applying parameters to interface
	var interfaceBytes []byte
	interfaceBytes, err = templates.ExecuteTemplate(tmpl, parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to execute interface file: %s\n%s", interfaceFile, err)
	}

	// Get values object from generated values yaml file
	var result = new(map[string]any)
	err = yaml.BytesToObject(interfaceBytes, result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated values from interface file: %s\n%s", interfaceFile, err)
	}

	return result, nil
}

// GetPackageFullNamesFromLocalRepository returns the list of package names in a local package repository, in alphabetical order.
func GetPackageFullNamesFromLocalRepository(repoDir string) ([]string, error) {
	var err error
	var ok bool

	var packagesDir = GetRepoPackagesDir(repoDir)

	// Exit early if the packages directory doesn't exist
	err = files.DirExists(packagesDir, "packages repository")
	if err != nil {
		return []string{}, err
	}

	// Get the items inside the repo directory.
	var namespaceDirEntries []fs.DirEntry
	namespaceDirEntries, err = os.ReadDir(packagesDir)
	if err != nil {
		return nil, err
	}

	// Traverse the packages directory.
	var packages = treeset.NewWithStringComparator()
	for _, namespaceDirEntry := range namespaceDirEntries {
		var namespaceDirAbsPath = filepath.Join(packagesDir, namespaceDirEntry.Name())

		// Skip files.
		if !namespaceDirEntry.IsDir() {
			log.Warningf("Found a file in the namespace directory '%s'", namespaceDirAbsPath)
			continue
		}

		var packageDirEntries []fs.DirEntry
		packageDirEntries, err = os.ReadDir(namespaceDirAbsPath)
		if err != nil {
			log.Warningf(
				"Failed to read the contents of namespace directory '%s': %s",
				namespaceDirAbsPath,
				err,
			)
			continue
		}

		for _, packageDirEntry := range packageDirEntries {
			var packageAbsPath = filepath.Join(namespaceDirAbsPath, packageDirEntry.Name())

			// Get the expected package name based on the path.
			var expectedPackageFullName string
			expectedPackageFullName, err = filepath.Rel(packagesDir, packageAbsPath)
			if err != nil {
				// This should never fail (we just calculated the path).
				log.Panicf(
					"Unexpectedly failed to get relative path: filepath.Rel(\"%s\", \"%s\")",
					namespaceDirAbsPath,
					packageAbsPath,
				)
			}
			// We always expect forward slashes.
			expectedPackageFullName = filepath.ToSlash(expectedPackageFullName)

			// Skip files.
			if !packageDirEntry.IsDir() {
				log.Warningf("Found a file in the packages directory: %s", packageAbsPath)
				continue
			}

			// Check if this is a valid package directory
			var packageDirAbsPath = GetPackageDir(repoDir, expectedPackageFullName)
			var packageInfo *PackageInfo
			packageInfo, err = GetPackageInfo(packageDirAbsPath)
			if err != nil {
				log.Warningf("Found invalid package '%s': %s", packageAbsPath, err)
				continue
			}

			// Get the package's full name from the package info file.
			var packageFullName = GetPackageFullName(packageInfo.Name, packageInfo.Version)

			// Check that the name of the directory matches the package's full name
			if expectedPackageFullName != packageFullName {
				// Log a warning
				log.Warningf(
					"Directory name '%s' does not match package name '%s': %s",
					packageDirEntry.Name(),
					packageFullName,
					packageAbsPath,
				)

				// Don't return this package, just continue looking for other packages
				continue
			}

			// Found a valid package, so add it to the list of found packages
			packages.Add(packageFullName)
		}
	}

	var result = make([]string, packages.Size())
	var i = 0
	for it := packages.Iterator(); it.Next(); {
		var packageName string
		var value = it.Value()
		packageName, ok = value.(string)
		if !ok {
			log.Panicf("Unexpected non-string value in tree map: %s", value)
		}

		result[i] = packageName
		i++
	}

	return result, nil
}
