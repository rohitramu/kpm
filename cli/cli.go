package cli

import (
	"fmt"

	"github.com/rohitramu/kpm/cli/implementation/cli_cobra"
	"github.com/rohitramu/kpm/cli/model/commands"
	"github.com/rohitramu/kpm/cli/model/utils/config"
)

func Execute() (err error) {
	var kpmConfig *config.KpmConfig
	kpmConfig, err = config.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read KPM configuration: %s", err)
	}

	var executeFunc = cli_cobra.GetCobraImplementation(kpmConfig, commands.Kpm)
	err = executeFunc()

	return err
}
