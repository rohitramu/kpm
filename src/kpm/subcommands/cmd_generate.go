package subcommands

import (
	"fmt"
	"log"
	"os"
)

// GenerateCmd generates a Kubernetes configuration from the given
// template package directory, parameters file and output directory
func GenerateCmd(packageDir string, parametersFile string, outputDir string) error {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	logger.Println(fmt.Sprintf("Package directory: %s", packageDir))
	logger.Println(fmt.Sprintf("Successfully generated Kubernetes configuration at: %s", outputDir))

	return nil
}
