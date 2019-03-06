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
	packageNameFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, pkg", constants.PackageNameFlagName),
		Usage: "Name of the package",
	}
	packageVersionFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, ver", constants.PackageVersionFlagName),
		Usage: "Version of the package",
	}
	packageDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, dir", constants.PackageDirFlagName),
		Usage: "Directory of the KPM package (defaults to current working directory)",
	}
	parametersFileFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, param", constants.ParametersFileFlagName),
		Usage: "Filepath of the parameters file to use",
	}
	outputNameFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, name", constants.OutputNameFlagName),
		Usage: "Name of the output configuration (defaults to <package name>-<package version)>",
	}
	outputDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, out", constants.OutputDirFlagName),
		Usage: "Directory in which output files should be written (defaults to current working directory) - WARNING: the <outputDir>/<outputName> directory will be deleted before generation",
	}
	kpmHomeDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, home", constants.KpmHomeDirFlagName),
		Usage: "Directory to use as the KPM home folder (defaults to ~/.kpm) - this value should be consistent across all invocations of KPM commands",
	}
	dockerRegistryURLFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, remote", constants.DockerRegistryURLFlagName),
		Usage: "The Docker registry URL to use when pulling or pushing a template package",
	}
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
		{
			Name:  constants.PackCmdName,
			Usage: "Packs a template package to make it available for use",
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
		{
			Name:    constants.GenerateCmdName,
			Aliases: []string{"gen"},
			Usage:   fmt.Sprintf("Generates a Kubernetes configuration using the template package specified by the \"--%s\" argument", constants.PackageDirFlagName),
			Flags: []cli.Flag{
				packageNameFlag,
				packageVersionFlag,
				parametersFileFlag,
				outputNameFlag,
				outputDirFlag,
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var packageName = getStringFlag(c, constants.PackageNameFlagName)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				var paramFile = getStringFlag(c, constants.ParametersFileFlagName)
				var outputName = getStringFlag(c, constants.OutputNameFlagName)
				var outputDir = getStringFlag(c, constants.OutputDirFlagName)
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				return subcommands.GenerateCmd(packageName, packageVersion, paramFile, outputName, outputDir, kpmHomeDir)
			},
		},
		{
			Name:  constants.PullCmdName,
			Usage: "Pulls a template package from a remote Docker registry",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
				dockerRegistryURLFlag,
				dockerNamespaceFlag,
				packageNameFlag,
				packageVersionFlag,
			},
			Action: func(c *cli.Context) error {
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var dockerRegistryURL = getStringFlag(c, constants.DockerRegistryURLFlagName)
				var dockerNamespace = getStringFlag(c, constants.DockerNamespaceFlagName)
				var packageName = getStringFlag(c, constants.PackageNameFlagName)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				return subcommands.PullCmd(kpmHomeDir, dockerRegistryURL, dockerNamespace, packageName, packageVersion)
			},
		},
		{
			Name:  constants.PushCmdName,
			Usage: "Pushes the template package to a remote Docker registry",
			Flags: []cli.Flag{
				kpmHomeDirFlag,
				dockerRegistryURLFlag,
				dockerNamespaceFlag,
				packageNameFlag,
				packageVersionFlag,
			},
			Action: func(c *cli.Context) error {
				var kpmHomeDir = getStringFlag(c, constants.KpmHomeDirFlagName)
				var dockerRegistryURL = getStringFlag(c, constants.DockerRegistryURLFlagName)
				var dockerNamespace = getStringFlag(c, constants.DockerNamespaceFlagName)
				var packageName = getStringFlag(c, constants.PackageNameFlagName)
				var packageVersion = getStringFlag(c, constants.PackageVersionFlagName)
				return subcommands.PushCmd(kpmHomeDir, dockerRegistryURL, dockerNamespace, packageName, packageVersion)
			},
		},
	}

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
