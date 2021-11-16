package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/rohitramu/kpm/subcommands"
	"github.com/rohitramu/kpm/subcommands/utils/constants"
	"github.com/rohitramu/kpm/subcommands/utils/log"
)

// Global flags
var (
	// Log level
	logLevelFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, log", constants.LogLevelFlagName),
		Usage: "The minimum severity log level to output - the severity log levels in increasing order are: \"debug\", \"verbose\", \"info\" (default), \"warning\", \"error\", \"none\"",
	}
)

// Flags
var (
	// Confirmation flag
	skipConfirmationFlag = cli.BoolFlag{
		Name:  fmt.Sprintf("%s, confirm", constants.SkipConfirmationFlagName),
		Usage: "Skip user confirmations.",
	}

	// Package version
	packageVersionFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, v", constants.PackageVersionFlagName),
		Usage: "Version of the package.",
	}

	// Parameters file
	parametersFileFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, f", constants.ParametersFileFlagName),
		Usage: "Filepath of the parameters file to use.",
	}

	// Output name
	outputNameFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, n", constants.OutputNameFlagName),
		Usage: "Name of the output (defaults to \"<package name>-<package version>\").",
	}

	// Output directory
	outputDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, d", constants.OutputDirFlagName),
		Usage: fmt.Sprintf("Directory in which output files should be written (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<outputName>\" will be deleted if it exists.", constants.GeneratedDirName),
		//TODO: Add support for environment variable and file config
	}

	// Export name
	exportNameFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, n", constants.ExportNameFlagName),
		Usage: "Name of the exported output (defaults to \"<package name>-<package version>\").",
	}

	// Export directory
	exportDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, d", constants.ExportDirFlagName),
		Usage: fmt.Sprintf("Directory in which exported files should be written (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<exportName>\" will be deleted if it exists.", constants.ExportDirName),
		//TODO: Add support for environment variable and file config
	}

	// KPM home directory
	kpmHomeDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, home", constants.KpmHomeDirFlagName),
		Usage: "Directory to use as the KPM home folder (defaults to \"~/.kpm\").",
		//TODO: Add support for environment variable and file config
	}

	// Docker registry
	dockerRegistryFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, r", constants.DockerRegistryFlagName),
		Usage: "The Docker registry to use when pulling or pushing a template package.",
		//TODO: Add support for environment variable and file config
	}
)

// Entrypoint
func main() {
	var err error

	// CLI app details
	app := cli.NewApp()
	app.Name = "kpm"
	app.Usage = "Kubernetes Package Manager"
	app.Version = "1.0.0"
	app.EnableBashCompletion = true

	// Global flags
	app.Flags = []cli.Flag{
		logLevelFlag,
	}

	// Sub-commands
	app.Commands = []cli.Command{
		// List
		{
			Name:    constants.ListCmdName,
			Aliases: []string{"ls"},
			Usage:   "Lists all packages that are currently available for use.",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				return subcommands.ListCmd(kpmHomeDir)
			},
		},

		// Remove
		{
			Name:    constants.RemoveCmdName,
			Aliases: []string{"rm"},
			Usage:   "Removes a package from the local KPM package repository.",
			Flags: []cli.Flag{
				packageVersionFlag,
				kpmHomeDirFlag,
				skipConfirmationFlag,
			},
			ArgsUsage: "<package name>",
			Action: func(c *cli.Context) error {
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var skipConfirmation = getBoolFlag(c, constants.SkipConfirmationFlagName)
				return subcommands.RemoveCmd(packageName, packageVersion, kpmHomeDir, skipConfirmation)
			},
		},

		// Purge
		{
			Name:  constants.PurgeCmdName,
			Usage: "Removes all versions of a package from the local KPM package repository.",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
				skipConfirmationFlag,
			},
			ArgsUsage: "<package name>",
			Action: func(c *cli.Context) error {
				var packageName = getStringArg(c, 0)
				var skipConfirmation = getBoolFlag(c, constants.SkipConfirmationFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				return subcommands.PurgeCmd(packageName, skipConfirmation, kpmHomeDir)
			},
		},

		// Pack
		{
			Name:  constants.PackCmdName,
			Usage: "Validates a template package and makes it available for use from the local KPM package repository.",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
				skipConfirmationFlag,
			},
			ArgsUsage: "<package directory>",
			Action: func(c *cli.Context) error {
				var packageDir = getStringArg(c, 0)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var skipConfirmation = getBoolFlag(c, constants.SkipConfirmationFlagName)
				return subcommands.PackCmd(packageDir, kpmHomeDir, skipConfirmation)
			},
		},

		// Unpack
		{
			Name:  constants.UnpackCmdName,
			Usage: "Exports a template package to the specified location.",
			Flags: []cli.Flag{
				packageVersionFlag,
				exportDirFlag,
				exportNameFlag,
				kpmHomeDirFlag,
				dockerRegistryFlag,
				skipConfirmationFlag,
			},
			ArgsUsage: "<package name>",
			Action: func(c *cli.Context) error {
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var exportDir = getStringFlag(c, constants.ExportDirFlagName)
				var exportName = getStringFlag(c, constants.ExportNameFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var dockerRegistry = getStringFlag(c, constants.DockerRegistryFlagName)
				var skipConfirmation = getBoolFlag(c, constants.SkipConfirmationFlagName)
				return subcommands.UnpackCmd(packageName, packageVersion, exportDir, exportName, kpmHomeDir, dockerRegistry, skipConfirmation)
			},
		},

		// Run
		{
			Name:  constants.RunCmdName,
			Usage: "Runs a template package.",
			Flags: []cli.Flag{
				packageVersionFlag,
				parametersFileFlag,
				outputDirFlag,
				outputNameFlag,
				kpmHomeDirFlag,
				dockerRegistryFlag,
				skipConfirmationFlag,
			},
			ArgsUsage: "<package name>",
			Action: func(c *cli.Context) error {
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var paramFile = getStringFlag(c, constants.ParametersFileFlagName)
				var outputDir = getStringFlag(c, constants.OutputDirFlagName)
				var outputName = getStringFlag(c, constants.OutputNameFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var dockerRegistry = getStringFlag(c, constants.DockerRegistryFlagName)
				var skipConfirmation = getBoolFlag(c, constants.SkipConfirmationFlagName)
				return subcommands.RunCmd(packageName, packageVersion, paramFile, outputDir, outputName, kpmHomeDir, dockerRegistry, skipConfirmation)
			},
		},

		// Inspect
		{
			Name:  constants.InspectCmdName,
			Usage: "Outputs the contents of the default parameters file in a template package.",
			Flags: []cli.Flag{
				packageVersionFlag,
				kpmHomeDirFlag,
				dockerRegistryFlag,
			},
			ArgsUsage: "<package name>",
			Action: func(c *cli.Context) error {
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var dockerRegistry = getStringFlag(c, constants.DockerRegistryFlagName)
				return subcommands.InspectCmd(packageName, packageVersion, kpmHomeDir, dockerRegistry)
			},
		},

		{
			Name:  constants.DockerCmdName,
			Usage: "Docker integration.",
			Subcommands: []cli.Command{
				// Push
				{
					Name:  constants.PushCmdName,
					Usage: "Pushes the template package to a remote Docker registry.",
					Flags: []cli.Flag{
						packageVersionFlag,
						kpmHomeDirFlag,
						dockerRegistryFlag,
					},
					ArgsUsage: "<package name>",
					Action: func(c *cli.Context) error {
						var packageName = getStringArg(c, 0)
						var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
						var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
						var dockerRegistry = getStringFlag(c, constants.DockerRegistryFlagName)
						return subcommands.PushCmd(packageName, packageVersion, kpmHomeDir, dockerRegistry)
					},
				},

				// Pull
				{
					Name:  constants.PullCmdName,
					Usage: "Pulls a template package from a remote Docker registry.",
					Flags: []cli.Flag{
						packageVersionFlag,
						kpmHomeDirFlag,
						dockerRegistryFlag,
					},
					ArgsUsage: "<package name>",
					Action: func(c *cli.Context) error {
						var packageName = getStringArg(c, 0)
						var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
						var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
						var dockerRegistry = getStringFlag(c, constants.DockerRegistryFlagName)
						return subcommands.PullCmd(packageName, packageVersion, kpmHomeDir, dockerRegistry)
					},
				},
			},
		},
	}

	// Do setup
	app.Before = func(c *cli.Context) error {
		// Set log level
		var logLevel = getGlobalStringFlag(c, constants.LogLevelFlagName)
		if logLevel != nil {
			// Parse the log level string
			var parsedLevel log.Level
			parsedLevel, err = log.Parse(*logLevel)
			if err != nil {
				return err
			}

			log.SetLevel(parsedLevel)
		}

		return nil
	}

	// Start the app
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal("%s", err)
	}
}

// +---------+
// | HELPERS |
// +---------+

func getGlobalStringFlag(c *cli.Context, flagName string) *string {
	if !c.GlobalIsSet(flagName) {
		return nil
	}

	var result = c.GlobalString(flagName)
	return &result
}

func getStringFlag(c *cli.Context, flagName string) *string {
	if !c.IsSet(flagName) {
		return nil
	}

	var result = c.String(flagName)
	return &result
}

func getStringArg(c *cli.Context, index int) *string {
	if index < 0 {
		log.Panic("Index cannot be negative: %d", index)
	}

	if c.NArg() <= index {
		return nil
	}

	var result = c.Args().Get(index)
	return &result
}

func getBoolFlag(c *cli.Context, flagName string) *bool {
	if !c.IsSet(flagName) {
		return nil
	}

	var result = c.Bool(flagName)
	return &result
}
