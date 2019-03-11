package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"./subcommands"
	"./subcommands/utils/constants"
	"./subcommands/utils/log"
)

// Flags
var (
	// Log level
	logLevelFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, log", constants.LogLevelFlagName),
		Usage: "The minimum severity log level to output - severities are in the following order: \"verbose\", \"info\" (default), \"warning\", \"error\", \"none\"",
	}

	// Package version
	packageVersionFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, v", constants.PackageVersionFlagName),
		Usage: "Version of the package",
	}

	// Parameters file
	parametersFileFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, f", constants.ParametersFileFlagName),
		Usage: "Filepath of the parameters file to use",
	}

	// Output name
	outputNameFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, n", constants.OutputNameFlagName),
		Usage: "Name of the output configuration (defaults to \"<package name>-<package version>\")",
	}

	// Output directory
	outputDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, d", constants.OutputDirFlagName),
		Usage: "Directory in which output files should be written (defaults to current working directory) - WARNING: the sub-directory specified by \"<outputName>\" will be deleted before generation if it exists",
	}

	// KPM home directory
	kpmHomeDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, home", constants.KpmHomeDirFlagName),
		Usage: "Directory to use as the KPM home folder (defaults to \"~/.kpm\")",
	}

	// Docker registry URL
	dockerRegistryFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, r", constants.DockerRegistryFlagName),
		Usage: "The Docker registry URL to use when pulling or pushing a template package",
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
			Usage:   "Lists all packages that are currently available for use",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				return subcommands.ListCmd(kpmHomeDir)
			},
		},

		// Pack
		{
			Name:  constants.PackCmdName,
			Usage: "Validates a template package and makes it available for use from the local KPM package repository",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var packageDir = getStringArg(c, 0)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				return subcommands.PackCmd(packageDir, kpmHomeDir)
			},
		},

		// Run
		{
			Name:  constants.RunCmdName,
			Usage: "Runs a template package",
			Flags: []cli.Flag{
				packageVersionFlag,
				parametersFileFlag,
				outputNameFlag,
				outputDirFlag,
				kpmHomeDirFlag,
				dockerRegistryFlag,
			},
			Action: func(c *cli.Context) error {
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var paramFile = getStringFlag(c, constants.ParametersFileFlagName)
				var outputName = getStringFlag(c, constants.OutputNameFlagName)
				var outputDir = getStringFlag(c, constants.OutputDirFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var dockerRegistry = getStringFlag(c, constants.DockerRegistryFlagName)
				return subcommands.RunCmd(packageName, packageVersion, paramFile, outputName, outputDir, kpmHomeDir, dockerRegistry)
			},
		},

		// Push
		{
			Name:  constants.PushCmdName,
			Usage: "Pushes the template package to a remote Docker registry",
			Flags: []cli.Flag{
				dockerRegistryFlag,
				packageVersionFlag,
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var dockerRegistry = getStringFlag(c, constants.DockerRegistryFlagName)
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				return subcommands.PushCmd(dockerRegistry, packageName, packageVersion, kpmHomeDir)
			},
		},

		// Pull
		{
			Name:  constants.PullCmdName,
			Usage: "Pulls a template package from a remote Docker registry",
			Flags: []cli.Flag{
				dockerRegistryFlag,
				packageVersionFlag,
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var dockerRegistry = getStringFlag(c, constants.DockerRegistryFlagName)
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				return subcommands.PullCmd(dockerRegistry, packageName, packageVersion, kpmHomeDir)
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
