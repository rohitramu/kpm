package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/rohitramu/kpm/pkg/utils/templates"
	"github.com/rohitramu/kpm/pkg/utils/validation"
)

// RunCmd runs the given template package directory and parameters file,
// and then writes the output files to the given output directory.
func RunCmd(packageName string, packageVersion string, optionalParametersFilePath *string, outputDirPath string, optionalOutputName *string, kpmHomeDirPath string, userHasConfirmed bool) error {
	var err error

	// Get KPM home directory
	var kpmHomeDir string
	kpmHomeDir, err = files.GetAbsolutePath(kpmHomeDirPath)
	if err != nil {
		return err
	}

	// Validate package name
	err = validation.ValidatePackageName(packageName)
	if err != nil {
		return err
	}

	// Validate package version
	err = validation.ValidatePackageVersion(packageVersion)
	if err != nil {
		return err
	}

	// Resolve generation paths
	var packageFullName = template_package.GetPackageFullName(packageName, packageVersion)
	var packageDirPath = template_package.GetPackageDir(kpmHomeDir, packageFullName)
	var outputName = validation.GetStringOrDefault(optionalOutputName, template_package.GetDefaultOutputName(packageName, packageVersion))
	var parametersFilePath string
	parametersFilePath, err = files.GetAbsolutePathOrDefault(optionalParametersFilePath, template_package.GetDefaultParametersFile(packageDirPath))
	if err != nil {
		return err
	}

	var packageOutputDirPath = filepath.Join(outputDirPath, outputName)

	// Log resolved values
	log.Verbosef("====")
	log.Verbosef("Package name:              %s", packageName)
	log.Verbosef("Package version:           %s", packageVersion)
	log.Verbosef("Package directory:         %s", packageDirPath)
	log.Verbosef("Parameters file:           %s", parametersFilePath)
	log.Verbosef("Output name:               %s", outputName)
	log.Verbosef("Output directory:          %s", outputDirPath)
	log.Verbosef("Package output directory:  %s", packageOutputDirPath)
	log.Verbosef("====")

	// Get the default parameters
	var packageParameters *templates.GenericMap
	packageParameters, err = template_package.GetPackageParameters(parametersFilePath)
	if err != nil {
		return err
	}

	// Get the dependency tree
	var dependencyTree *template_package.DependencyTree
	if dependencyTree, err = template_package.GetDependencyTree(kpmHomeDir, packageName, packageVersion, outputName, packageParameters); err != nil {
		return err
	}

	// Delete the output directory in case it isn't empty
	if err = files.DeleteDirIfExists(packageOutputDirPath, "output", userHasConfirmed); err != nil {
		return err
	}

	// Execute template packages in the dependency tree and write the output to the filesystem
	var numPackages int
	numPackages, err = dependencyTree.VisitNodesDepthFirst(func(relativeFilePath []string, friendlyNamePath []string, executableTemplates []*template.Template, templateInput *templates.GenericMap) error {
		// Get the output directory
		var outputDir = filepath.Join(outputDirPath, filepath.Join(relativeFilePath...))

		// Create the output directory if it doesn't exist
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return err
		}

		// Get the templates in the package
		for _, tmpl := range executableTemplates {
			// Execute the template with the provided input data
			var templateOutput []byte
			templateOutput, err = templates.ExecuteTemplate(tmpl, templateInput)
			if err != nil {
				return fmt.Errorf("failed to execute package: %s\n%s", strings.Join(friendlyNamePath, " -> "), err)
			}

			// Write the data to the filesystem
			var outputFilePath = filepath.Join(outputDir, tmpl.Name())
			log.Verbosef("Writing file: %s", outputFilePath)
			os.WriteFile(outputFilePath, templateOutput, 0755)
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Debugf("Executed %d packages", numPackages)

	return nil
}
