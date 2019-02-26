package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/urfave/cli"

	"./subcommands"
)

// Main logger
var logger = log.New(os.Stderr, "", log.LstdFlags)

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
				parametersFileFlag,
				outputDirFlag,
			},
			Action: func(c *cli.Context) error {
				packageDir := getFlagValuePackageDir(c)
				paramFile := getFlagValueParametersFile(c, packageDir)
				outputDir := getFlagValueOutputDir(c, packageDir)
				return subcommands.GenerateCmd(packageDir, paramFile, outputDir)
			},
		},
		{
			Name:  "push",
			Usage: "Pushes the template package to a docker repository",
			Flags: []cli.Flag{
				packageDirFlag,
			},
			Action: func(c *cli.Context) error {
				packageDir := getFlagValuePackageDir(c)
				return subcommands.PushCmd(packageDir)
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
				packageDir := getFlagValuePackageDir(c)
				return subcommands.ApplyCmd(packageDir)
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

// Value for flag "packageDir"
func getFlagValuePackageDir(c *cli.Context) *string {
	var err error

	// Get user input
	var packageDir = getStringFlag(c, &packageDirFlagName)

	// Get default path (current working directory)
	defaultPath, err := os.Getwd()
	if err != nil {
		logger.Fatalln(err)
	}

	// Resolve absolute path to use for packageDir
	packageDir = getAbsolutePathOrDefaultOrExit(packageDir, &defaultPath)

	return packageDir
}

// Value for flag "paramFile"
func getFlagValueParametersFile(c *cli.Context, defaultDir *string) *string {
	var defaultPath = filepath.Join(*defaultDir, "parameters.yaml")

	var paramFile = getStringFlag(c, &parametersFileFlagName)
	paramFile = getAbsolutePathOrDefaultOrExit(paramFile, &defaultPath)

	return paramFile
}

// Value for flag "outputDir"
func getFlagValueOutputDir(c *cli.Context, defaultDir *string) *string {
	var defaultPath = filepath.Join(*defaultDir, "_output_")

	var outputDir = getStringFlag(c, &outputDirFlagName)
	outputDir = getAbsolutePathOrDefaultOrExit(outputDir, &defaultPath)

	return outputDir
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

func getAbsolutePathOrDefaultOrExit(path *string, defaultPath *string) *string {
	var outputPath *string
	if path != nil {
		outputPath = getAbsolutePathOrExit(path)
	} else {
		outputPath = getAbsolutePathOrExit(defaultPath)
	}

	return outputPath
}

func getAbsolutePathOrExit(path *string) *string {
	var err error

	var outputPath = path

	// Resolve "~" to the user's home directory if required
	if len(*outputPath) > 0 && (*outputPath)[0] == '~' && ((*outputPath)[1] == '/' || (*outputPath)[1] == '\\') {
		var usr *(user.User)
		usr, err = user.Current()
		*outputPath = filepath.Join(usr.HomeDir, (*outputPath)[2:])
	}

	// Check if path is already absolute
	if !filepath.IsAbs(*outputPath) {
		// Get absolute path
		var newOutputPath string
		newOutputPath, err = filepath.Abs(*outputPath)
		outputPath = &newOutputPath

		// Exit on error
		if err != nil {
			logger.Fatalln(err)
		}
	}

	return outputPath
}
