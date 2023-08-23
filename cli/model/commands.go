package model

import (
	"errors"
	"fmt"

	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
)

var KpmCmd = &Command{
	Name: "kpm",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			logLevelFlag,
		},
	},
	SubCommands: []*Command{
		listCmd,
		removeCmd,
		purgeCmd,
		packCmd,
		unpackCmd,
		inspectCmd,
		runCmd,
		newCmd,
		repoCmd,
	},
}

var listCmd = &Command{
	Name:             "list",
	Alias:            "ls",
	ShortDescription: "Lists all template packages.",
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Get the list of package names.
		var packages []string
		packages, err = pkg.ListCmd(kpmHomeDir)
		if err != nil {
			return err
		}

		// Print package names.
		for _, packageName := range packages {
			log.Outputf(packageName)
		}

		return nil
	},
}

var removeCmd = &Command{
	Name:             "remove",
	Alias:            "rm",
	ShortDescription: "Removes a template package.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{packageVersionFlag},
		BoolFlags:   []Flag[bool]{userConfirmationFlag},
	},
	Args: ArgCollection{
		MandatoryArgs: []*Arg{{
			Name:             "package-name",
			ShortDescription: "The name of the template package to remove.",
		}},
	},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var shouldSkipUserConfirmation = userConfirmationFlag.GetValueOrDefault(config)
		var packageVersion = packageVersionFlag.GetValueOrDefault(config)

		// Args
		var packageName string = args.MandatoryArgs[0].Value

		if packageVersion == "" {
			// Since the package version was not provided, check the local repository for the highest version.
			var err error
			if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
				return fmt.Errorf("could not find package '%s' in the local KPM repository: %s", packageName, err)
			}
		}

		return pkg.RemoveCmd(
			packageName,
			packageVersion,
			kpmHomeDir,
			shouldSkipUserConfirmation)
	},
}

// TODO: Merge the "purge" command into the "remove" command (use flags to determine behavior).
var purgeCmd = &Command{
	Name:             "purge",
	ShortDescription: "Removes all versions of a template package.",
	Flags: FlagCollection{
		BoolFlags: []Flag[bool]{userConfirmationFlag},
	},
	Args: ArgCollection{
		OptionalArg: &Arg{
			Name:             "package-name",
			ShortDescription: "The name of the template package to purge.  If this is not provided, all versions of all template packages will be deleted.",
		},
	},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var skipConfirmation = userConfirmationFlag.GetValueOrDefault(config)

		// Args
		var packageName string
		if args.OptionalArg != nil {
			packageName = args.OptionalArg.Value
		}

		return pkg.PurgeCmd(packageName, skipConfirmation, kpmHomeDir)
	},
}

var packCmd = &Command{
	Name:             "pack",
	ShortDescription: "Validates a template package and makes it available for use.",
	Flags: FlagCollection{
		BoolFlags: []Flag[bool]{userConfirmationFlag},
	},
	Args: ArgCollection{
		MandatoryArgs: []*Arg{{
			Name:             "package-directory",
			ShortDescription: "The location of the template package directory which should be packed.",
		}},
	},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var skipConfirmation = userConfirmationFlag.GetValueOrDefault(config)

		// Args
		var packageDir = args.MandatoryArgs[0].Value

		return pkg.PackCmd(packageDir, kpmHomeDir, skipConfirmation)
	},
}

var unpackCmd = &Command{
	Name:             "unpack",
	ShortDescription: "Exports a template package to the specified location.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			packageVersionFlag,
			exportDirFlag,
			exportNameFlag,
		},
		BoolFlags: []Flag[bool]{
			userConfirmationFlag,
		},
	},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var skipConfirmation = userConfirmationFlag.GetValueOrDefault(config)
		var packageVersion = packageVersionFlag.GetValueOrDefault(config)
		var exportDir = exportDirFlag.GetValueOrDefault(config)
		var exportName = exportNameFlag.GetValueOrDefault(config)

		// Args
		var packageName = args.MandatoryArgs[0].Value

		// Validation
		{
			// Package version
			if packageVersion == "" {
				// Since the package version was not provided, check the local repository for the highest version
				var err error
				if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
					return fmt.Errorf("package version must be provided if the package does not exist in the local repository: %s", err)
				}
			}

			// Export name
			if exportName == "" {
				exportName = template_package.GetDefaultExportName(packageName, packageVersion)
			}
		}

		return pkg.UnpackCmd(packageName, packageVersion, exportDir, exportName, kpmHomeDir, skipConfirmation)
	},
}

var inspectCmd = &Command{
	Name:             "inpect",
	ShortDescription: "Prints the contents of the default parameters file in a template package.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			packageVersionFlag,
		},
		BoolFlags: []Flag[bool]{
			userConfirmationFlag,
		},
	},
	Args: ArgCollection{
		MandatoryArgs: []*Arg{{
			Name:             "package-name",
			ShortDescription: "The name of the template package to run.",
		}},
	},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Args
		var packageName = args.MandatoryArgs[0].Value

		// Flags
		var packageVersion = packageVersionFlag.GetValueOrDefault(config)

		// Validation
		{
			// Package version
			if packageVersion == "" {
				// Since the package version was not provided, check the local repository for the highest version.
				var err error
				if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
					return fmt.Errorf("could not find package '%s' in the local KPM repository: %s", packageName, err)
				}
			}
		}

		return pkg.InspectCmd(packageName, packageVersion, kpmHomeDir)
	},
}

var runCmd = &Command{
	Name:             "run",
	ShortDescription: "Runs a template package.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			packageVersionFlag,
			parametersFileFlag,
			outputDirFlag,
			outputNameFlag,
		},
		BoolFlags: []Flag[bool]{
			userConfirmationFlag,
		},
	},
	Args: ArgCollection{
		MandatoryArgs: []*Arg{{
			Name:             "package-name",
			ShortDescription: "The name of the template package to run.",
		}},
	},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Args
		var packageName = args.MandatoryArgs[0].Value

		// Flags
		var packageVersion = packageVersionFlag.GetValueOrDefault(config)
		var paramFile = parametersFileFlag.GetValueOrDefault(config)
		var outputDir = outputDirFlag.GetValueOrDefault(config)
		var outputName = outputNameFlag.GetValueOrDefault(config)
		var skipConfirmation = userConfirmationFlag.GetValueOrDefault(config)

		// Validation
		var optionalParamFile = &paramFile
		var optionalOutputName = &outputName
		{
			// Package version
			if packageVersion == "" {
				// Since the package version was not provided, check the local repository for the highest version.
				var err error
				if packageVersion, err = template_package.GetHighestPackageVersion(kpmHomeDir, packageName); err != nil {
					return fmt.Errorf("could not find package '%s' in the local KPM repository: %s", packageName, err)
				}
			}
			// Parameters file
			if paramFile == "" {
				optionalParamFile = nil
			}
			// Output name
			if outputName == "" {
				optionalOutputName = nil
			}
		}

		return pkg.RunCmd(packageName, packageVersion, optionalParamFile, outputDir, optionalOutputName, kpmHomeDir, skipConfirmation)
	},
}

var newCmd = &Command{
	Name:             "new-package",
	Alias:            "new",
	ShortDescription: "Creates a new template package.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			newPackageOutputDirFlag,
		},
		BoolFlags: []Flag[bool]{
			userConfirmationFlag,
		},
	},
	Args: ArgCollection{
		OptionalArg: &Arg{
			Name:             "package-name",
			ShortDescription: "The name of the new package.",
			Value:            "hello-kpm",
		},
	},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) error {
		// Flags
		var skipConfirmation = userConfirmationFlag.GetValueOrDefault(config)
		var packageDir = newPackageOutputDirFlag.GetValueOrDefault(config)

		// Args
		var packageName = args.OptionalArg.Value

		return pkg.NewTemplatePackageCmd(packageName, packageDir, skipConfirmation)
	},
}

var repoCmdName = "repositories"
var repoCmd = &Command{
	Name:             repoCmdName,
	Alias:            "repo",
	ShortDescription: "Commands for interacting with template package repositories.",
	SubCommands: []*Command{
		repoListCmd,
		repoPushCmd,
	},
}

var repoListCmd = &Command{
	Name:             "list",
	Alias:            "ls",
	ShortDescription: "Lists the names of available repositories.",
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var repos []string
		repos, err = pkg.ListPackageRepositories(&config.Repositories)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			log.Outputf(repo)
		}

		return nil
	},
}

var repoPushCmd = &Command{
	Name:             "push",
	ShortDescription: "Pushes a template package to a repository.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			repoNameFlag,
			packageVersionFlag,
		},
		BoolFlags: []Flag[bool]{
			userConfirmationFlag,
		},
	},
	Args: ArgCollection{MandatoryArgs: []*Arg{
		{
			Name:             "package-name",
			ShortDescription: "The name of the template package to push.",
			IsValidFunc:      validatePackageName,
		},
	}},
	ExecuteFunc: func(config *KpmConfig, args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		var packageName = args.MandatoryArgs[0].Value
		var packageVersion = packageVersionFlag.GetValueOrDefault(config)

		// If the repo name isn't provided, pick the first one.
		var repoName = repoNameFlag.GetValueOrDefault(config)
		if repoName == "" {
			var repoNames = config.Repositories.GetRepositoryNames()
			if len(repoNames) == 0 {
				return errors.New("no repositories configured")
			}

			repoName = repoNames[0]
		}

		return pkg.PushToRepository(
			kpmHomeDir,
			&config.Repositories,
			repoName,
			packageName,
			packageVersion,
			userConfirmationFlag.GetValueOrDefault(config),
		)
	},
}
