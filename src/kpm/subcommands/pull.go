package subcommands

import (
	"fmt"
	"log"
	"os"
)

// PullCmd pulls a template package from a Docker registry to the local filesystem
func PullCmd(packageDir *string) error {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	logger.Println(fmt.Sprintf("Package directory: %s", *packageDir))
	return nil
}
