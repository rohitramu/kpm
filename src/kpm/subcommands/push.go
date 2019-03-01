package subcommands

import (
	"fmt"
	"log"
	"os"
)

// PushCmd pushes the template package to a Docker registry.
func PushCmd(packageName *string, kpmHomeDir *string) error {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	logger.Println(fmt.Sprintf("Package name: %s", *packageName))
	return nil
}
