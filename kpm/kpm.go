package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"./subcommands"
	"./subcommands/utils/constants"
	"./subcommands/utils/logger"
)

// Flags
var (
	// Package version
	packageVersionFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, ver", constants.PackageVersionFlagName),
		Usage: "Version of the package",
	}

	// Package directory
	packageDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, dir", constants.PackageDirFlagName),
		Usage: "Directory of the KPM package (defaults to current working directory)",
	}

	// Parameters file
	parametersFileFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, param", constants.ParametersFileFlagName),
		Usage: "Filepath of the parameters file to use",
	}

	// Output name
	outputNameFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, name", constants.OutputNameFlagName),
		Usage: "Name of the output configuration (defaults to <package name>-<package version)>",
	}

	// Output directory
	outputDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, out", constants.OutputDirFlagName),
		Usage: "Directory in which output files should be written (defaults to current working directory) - WARNING: the sub-directory <outputName> will be deleted before generation if it exists",
	}

	// KPM home directory
	kpmHomeDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, home", constants.KpmHomeDirFlagName),
		Usage: "Directory to use as the KPM home folder (defaults to ~/.kpm)",
	}

	// Docker registry URL
	dockerRegistryURLFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, remote", constants.DockerRegistryURLFlagName),
		Usage: "The Docker registry URL to use when pulling or pushing a template package",
	}

	// Docker namespace
	dockerNamespaceFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, ns", constants.DockerNamespaceFlagName),
		Usage: "The Docker namespace to use when pulling or pushing a template package",
	}
)

// Entrypoint
func main() {
	// CLI app details
	app := cli.NewApp()
	app.Name = "kpm"
	app.Usage = "Kubernetes Package Manager"
	app.Version = "1.0.0"

	// Sub-commands
	app.Commands = []cli.Command{
		// List
		{
			Name:  constants.ListCmdName,
			Usage: "Lists all packages that are currently available for use",
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
				packageDirFlag,
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var packageDir = getStringFlag(c, constants.PackageDirFlagName)
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
			},
			Action: func(c *cli.Context) error {
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var paramFile = getStringFlag(c, constants.ParametersFileFlagName)
				var outputName = getStringFlag(c, constants.OutputNameFlagName)
				var outputDir = getStringFlag(c, constants.OutputDirFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				return subcommands.RunCmd(packageName, packageVersion, paramFile, outputName, outputDir, kpmHomeDir)
			},
		},

		// Pull
		{
			Name:  constants.PullCmdName,
			Usage: "Pulls a template package from a remote Docker registry",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
				dockerRegistryURLFlag,
				dockerNamespaceFlag,
				packageVersionFlag,
			},
			Action: func(c *cli.Context) error {
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var dockerRegistryURL = getStringFlag(c, constants.DockerRegistryURLFlagName)
				var dockerNamespace = getStringFlag(c, constants.DockerNamespaceFlagName)
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				return subcommands.PullCmd(kpmHomeDir, dockerRegistryURL, dockerNamespace, packageName, packageVersion)
			},
		},

		// Push
		{
			Name:  constants.PushCmdName,
			Usage: "Pushes the template package to a remote Docker registry",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
				dockerRegistryURLFlag,
				dockerNamespaceFlag,
				packageVersionFlag,
			},
			Action: func(c *cli.Context) error {
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var dockerRegistryURL = getStringFlag(c, constants.DockerRegistryURLFlagName)
				var dockerNamespace = getStringFlag(c, constants.DockerNamespaceFlagName)
				var packageName = getStringArg(c, 0)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				return subcommands.PushCmd(kpmHomeDir, dockerRegistryURL, dockerNamespace, packageName, packageVersion)
			},
		},
	}

	// Start the app
	err := app.Run(os.Args)
	if err != nil {
		logger.Default.Error.Fatalln(err)
	}
}

// +---------+
// | HELPERS |
// +---------+

func getStringFlag(c *cli.Context, flagName string) *string {
	if !c.IsSet(flagName) {
		return nil
	}

	var result = c.String(flagName)
	return &result
}

func getStringArg(c *cli.Context, index int) *string {
	if index < 0 {
		logger.Default.Error.Panicln(fmt.Sprintf("Index cannot be negative: %d", index))
	}

	if c.NArg() <= index {
		return nil
	}

	var result = c.Args().Get(index)
	return &result
}
