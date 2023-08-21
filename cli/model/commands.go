package model

import (
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
	},
}

var listCmd = &Command{
	Name:             "list",
	Alias:            "ls",
	ShortDescription: "Lists all template packages.",
	IsValidFunc: func(args ArgCollection) (bool, error) {
		log.Infof("Command validation executed!")
		return true, nil
	},
	ExecuteFunc: func(args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		return pkg.ListCmd(kpmHomeDir)
	},
}

var removeCmd = &Command{
	Name:             "remove",
	Alias:            "rm",
	ShortDescription: "Removes a template package.",
	Flags: FlagCollection{
		StringFlags: []Flag[string]{packageVersionFlag},
		BoolFlags:   []Flag[bool]{skipUserConfirmationFlag},
	},
	Args: ArgCollection{
		MandatoryArgs: []*Arg{{
			Name:             "package-name",
			ShortDescription: "The name of the template package to remove.",
		}},
	},
	ExecuteFunc: func(args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var shouldSkipUserConfirmation = skipUserConfirmationFlag.GetValueOrDefault()
		var packageVersion = packageVersionFlag.GetValueOrDefault()

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
		BoolFlags: []Flag[bool]{skipUserConfirmationFlag},
	},
	Args: ArgCollection{
		OptionalArg: &Arg{
			Name:             "package-name",
			ShortDescription: "The name of the template package to purge.  If this is not provided, all versions of all template packages will be deleted.",
		},
	},
	ExecuteFunc: func(args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var skipConfirmation = skipUserConfirmationFlag.GetValueOrDefault()

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
		BoolFlags: []Flag[bool]{skipUserConfirmationFlag},
	},
	Args: ArgCollection{
		MandatoryArgs: []*Arg{{
			Name:             "package-directory",
			ShortDescription: "The location of the template package directory which should be packed.",
		}},
	},
	ExecuteFunc: func(args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var skipConfirmation = skipUserConfirmationFlag.GetValueOrDefault()

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
			skipUserConfirmationFlag,
		},
	},
	ExecuteFunc: func(args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Flags
		var skipConfirmation = skipUserConfirmationFlag.GetValueOrDefault()
		var packageVersion = packageVersionFlag.GetValueOrDefault()
		var exportDir = exportDirFlag.GetValueOrDefault()
		var exportName = exportNameFlag.GetValueOrDefault()

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
			skipUserConfirmationFlag,
		},
	},
	Args: ArgCollection{
		MandatoryArgs: []*Arg{{
			Name:             "package-name",
			ShortDescription: "The name of the template package to run.",
		}},
	},
	ExecuteFunc: func(args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Args
		var packageName = args.MandatoryArgs[0].Value

		// Flags
		var packageVersion = packageVersionFlag.GetValueOrDefault()

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
			skipUserConfirmationFlag,
		},
	},
	Args: ArgCollection{
		MandatoryArgs: []*Arg{{
			Name:             "package-name",
			ShortDescription: "The name of the template package to run.",
		}},
	},
	ExecuteFunc: func(args ArgCollection) (err error) {
		var kpmHomeDir string
		if kpmHomeDir, err = template_package.GetKpmHomeDir(); err != nil {
			return err
		}

		// Args
		var packageName = args.MandatoryArgs[0].Value

		// Flags
		var packageVersion = packageVersionFlag.GetValueOrDefault()
		var paramFile = parametersFileFlag.GetValueOrDefault()
		var outputDir = outputDirFlag.GetValueOrDefault()
		var outputName = outputNameFlag.GetValueOrDefault()
		var skipConfirmation = skipUserConfirmationFlag.GetValueOrDefault()

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
			skipUserConfirmationFlag,
		},
	},
	Args: ArgCollection{
		OptionalArg: &Arg{
			Name:             "package-name",
			ShortDescription: "The name of the new package.",
			Value:            "hello-kpm",
		},
	},
	ExecuteFunc: func(args ArgCollection) error {
		// Flags
		var skipConfirmation = skipUserConfirmationFlag.GetValueOrDefault()
		var packageDir = newPackageOutputDirFlag.GetValueOrDefault()

		// Args
		var packageName = args.OptionalArg.Value

		return pkg.NewTemplatePackageCmd(packageName, packageDir, skipConfirmation)
	},
}
