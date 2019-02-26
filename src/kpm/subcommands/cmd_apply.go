package subcommands

import (
	"fmt"
	"log"
	"os"
)

// ApplyCmd applies the generated Kubernetes configuration to a
// Kubernetes cluster
func ApplyCmd(packageDir *string) error {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	logger.Println(fmt.Sprintf("Package directory: %s", *packageDir))
	return nil
}
