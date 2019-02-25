package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	"./subcommands"
)

// Main logger
var logger = log.New(os.Stderr, "kpm: ", log.Ltime|log.Lshortfile)

// Sub-command names
var (
	generateCmdName = "generate"
)

// Flag names
var (
	packageDirFlagName     = "packageDir"
	parametersFileFlagName = "parametersFile"
	outputDirFlagName      = "outputDir"
)

// Flags
var (
	packageDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, p", packageDirFlagName),
		Usage: "Directory of the KPM package (defaults to current working directory)",
	}
	parametersFileFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, f", parametersFileFlagName),
		Usage: "Filepath of the parameters file to use",
	}
	outputDirFlag = cli.StringFlag{
		Name:  fmt.Sprintf("%s, o", outputDirFlagName),
		Usage: "Directory in which output files should be created",
	}
)

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
				parametersFileFlag,
				outputDirFlag,
			},
			Action: func(c *cli.Context) error {
				return subcommands.GenerateCmd(getFlagValuePackageDir(c), getFlagValueParamFile(c), getFlagValueOutputDir(c))
			},
		},
		{
			Name:  "push",
			Usage: "Pushes the template package to a docker repository",
			Flags: []cli.Flag{
				packageDirFlag,
			},
			Action: func(c *cli.Context) error {
				return subcommands.PushCmd(getFlagValuePackageDir(c))
			},
		},
		{
			Name:  "apply",
			Usage: "Applies a generated Kubernetes configuration to a Kubernetes cluster",
			Flags: []cli.Flag{
				packageDirFlag,
				parametersFileFlag,
			},
			Action: func(c *cli.Context) error {
				return subcommands.ApplyCmd(getFlagValuePackageDir(c))
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalln(err)
	}
}

// +-------+
// | FLAGS |
// +-------+

func getFlagValuePackageDir(c *cli.Context) string {
	var err error

	// Get user input
	packageDir := c.String(packageDirFlagName)

	// Get default path (current working directory)
	defaultPath, err := os.Getwd()
	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}

	// Resolve absolute path to use for packageDir
	packageDir = getAbsolutePathOrDefaultOrExit(packageDir, defaultPath)

	return packageDir
}

func getFlagValueParamFile(c *cli.Context) string {
	paramFile := c.String(parametersFileFlagName)
	paramFile = getAbsolutePathOrDefaultOrExit(paramFile, "")
	c.Set(parametersFileFlagName, paramFile)

	return paramFile
}

func getFlagValueOutputDir(c *cli.Context) string {
	outputDir := c.String(outputDirFlagName)
	outputDir = getAbsolutePathOrDefaultOrExit(outputDir, "../_generated_")
	c.Set(outputDirFlagName, outputDir)

	return outputDir
}

// +---------+
// | HELPERS |
// +---------+

func getAbsolutePathOrDefaultOrExit(path string, defaultPath string) string {
	if path != "" {
		path = getAbsolutePathOrExit(path)
	} else {
		path = getAbsolutePathOrExit(defaultPath)
	}

	return path
}

func getAbsolutePathOrExit(path string) string {
	var err error

	// Check if path is already absolute
	if !filepath.IsAbs(path) {
		// Get absolute path
		path, err = filepath.Abs(path)

		// Exit on error
		if err != nil {
			logger.Fatalln(err)
			os.Exit(1)
		}
	}

	return path
}
