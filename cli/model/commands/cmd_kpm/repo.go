package cmd_kpm

import (
	"github.com/rohitramu/kpm/cli/model/commands/cmd_kpm/cmd_kpm_repo"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/types"
)

var Repo = &types.Command{
	Name:             constants.CmdRepo,
	Alias:            "repo",
	ShortDescription: "Commands for interacting with template package repositories.",
	SubCommands: []*types.Command{
		cmd_kpm_repo.ListCmd,
		cmd_kpm_repo.FindCmd,
		cmd_kpm_repo.PushCmd,
		cmd_kpm_repo.PullCmd,
	},
}
