package subcommands

import (
	"fmt"

	"./utils/logger"
)

// PullCmd pulls a template package from a Docker registry to the local filesystem.
func PullCmd(kpmHomeDirArg *string, dockerRegistryURLArg *string, dockerNamespaceArg *string, packageNameArg *string, packageVersionArg *string) error {
	var err error

	logger.Default.Verbose.Println(fmt.Sprintf("Package name: %s", *packageNameArg))
	return err
}
