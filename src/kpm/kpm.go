package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"./subcommands"
)

// Main logger
var logger = log.New(os.Stdout, "", log.LstdFlags)

// Sub-command names
var (
	generateCmdName = "generate"
)

// Flag names
var (
	packageDirFlagName  = "packageDir"
	valuesFileFlagName  = "valuesFile"
	outputDirFlagName   = "outputDir"
	downloadDirFlagName = "downloadDir"
)

// Flags
var (
	packageDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, p", packageDirFlagName),
		Usage: "Directory of the KPM package (defaults to current working directory)",
	}
	valuesFileFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, v", valuesFileFlagName),
		Usage: "Filepath of the values file to use",
	}
	outputDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, o", outputDirFlagName),
		Usage: "Directory in which output files should be created",
	}
	downloadDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, d", downloadDirFlagName),
		Usage: "Directory in which KPM packages should be downloaded to (defaults to the current working directory)",
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
			Name:  generateCmdName,
			Usage: fmt.Sprintf("Generates a Kubernetes configuration using the template package specified by the \"--%s\" argument", packageDirFlagName),
			Flags: []cli.Flag{
				packageDirFlag,
				valuesFileFlag,
				outputDirFlag,
			},
			Action: func(c *cli.Context) error {
				var packageDir = getStringFlag(c, &packageDirFlagName)
				var paramFile = getStringFlag(c, &valuesFileFlagName)
				var outputDir = getStringFlag(c, &outputDirFlagName)
				return subcommands.GenerateCmd(packageDir, paramFile, outputDir)
			},
		},
		{
			Name:  "pull",
			Usage: "Pulls a template package from a docker repository",
			Flags: []cli.Flag{
				outputDirFlag,
			},
			Action: func(c *cli.Context) error {
				downloadDir := getStringFlag(c, &downloadDirFlagName)
				return subcommands.PullCmd(downloadDir)
			},
		},
		{
			Name:  "push",
			Usage: "Pushes the template package to a docker repository",
			Flags: []cli.Flag{
				packageDirFlag,
			},
			Action: func(c *cli.Context) error {
				packageDir := getStringFlag(c, &packageDirFlagName)
				return subcommands.PushCmd(packageDir)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalln(err)
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
