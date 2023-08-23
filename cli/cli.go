package cli

import (
	"fmt"

	"github.com/rohitramu/kpm/cli/implementation/cli_cobra"
	"github.com/rohitramu/kpm/cli/model"
)

func Execute() (err error) {
	var config *model.KpmConfig
	config, err = model.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read KPM configuration: %s", err)
	}

	var executeFunc = cli_cobra.GetCobraImplementation(config, model.KpmCmd)
	err = executeFunc()

	return err
}
