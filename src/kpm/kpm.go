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
		Name:  fmt.Sprintf("%s, n", constants.PackageNameFlagName),
		Usage: "Name of the package",
	}
	packageVersionFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, v", constants.PackageVersionFlagName),
		Usage: "Version of the package",
	}
	packageDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, d", constants.PackageDirFlagName),
		Usage: "Directory of the KPM package (defaults to current working directory)",
	}
	parametersFileFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, f", constants.ParametersFileFlagName),
		Usage: "Filepath of the parameters file to use",
	}
	outputDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, o", constants.OutputDirFlagName),
		Usage: "Directory in which output files should be written",
	}
	kpmHomeDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, k", constants.KpmHomeDirFlagName),
		Usage: "Directory to use as the KPM home folder - this value should be consistent value across all invocations of KPM commands",
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
				var kpmHomeDir = getStringFlag(c, &constants.KpmHomeDirFlagName)
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
				var packageDir = getStringFlag(c, &constants.PackageDirFlagName)
				var kpmHomeDir = getStringFlag(c, &constants.KpmHomeDirFlagName)
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
				outputDirFlag,
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var packageName = getStringFlag(c, &constants.PackageNameFlagName)
				var packageVersion = getStringFlag(c, &constants.PackageVersionFlagName)
				var paramFile = getStringFlag(c, &constants.ParametersFileFlagName)
				var outputDir = getStringFlag(c, &constants.OutputDirFlagName)
				var kpmHomeDir = getStringFlag(c, &constants.KpmHomeDirFlagName)
				return subcommands.GenerateCmd(packageName, packageVersion, paramFile, outputDir, kpmHomeDir)
			},
		},
		{
			Name:  constants.PullCmdName,
			Usage: "Pulls a template package from a docker repository",
			Flags: []cli.Flag{
				packageNameFlag,
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var packageName = getStringFlag(c, &constants.PackageNameFlagName)
				var kpmHomeDir = getStringFlag(c, &constants.KpmHomeDirFlagName)
				return subcommands.PullCmd(packageName, kpmHomeDir)
			},
		},
		{
			Name:  constants.PushCmdName,
			Usage: "Pushes the template package to a docker repository",
			Flags: []cli.Flag{
				packageNameFlag,
				kpmHomeDirFlag,
			},
			Action: func(c *cli.Context) error {
				var packageName = getStringFlag(c, &constants.PackageNameFlagName)
				var kpmHomeDir = getStringFlag(c, &constants.KpmHomeDirFlagName)
				return subcommands.PushCmd(packageName, kpmHomeDir)
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

func getStringFlag(c *cli.Context, flagName *string) *string {
	if !c.IsSet(*flagName) {
		return nil
	}

	var result = c.String(*flagName)
	return &result
}
