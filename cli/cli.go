package cli

import (
	"github.com/rohitramu/kpm/cli/implementation/cli_cobra"
	"github.com/rohitramu/kpm/cli/model"
)

func Execute() error {
	var executeFunc = cli_cobra.GetCobraImplementation(model.KpmCmd)
	return executeFunc()
}
