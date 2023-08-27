package commands

import (
	"github.com/rohitramu/kpm/cli/model/commands/cmd_kpm"
	"github.com/rohitramu/kpm/cli/model/flags"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/types"
)

var Kpm = &types.Command{
	Name: constants.CmdKpm,
	Flags: types.FlagCollection{
		StringFlags: []types.Flag[string]{
			flags.LogLevel,
		},
	},
	SubCommands: []*types.Command{
		cmd_kpm.List,
		cmd_kpm.Remove,
		cmd_kpm.Purge,
		cmd_kpm.Pack,
		cmd_kpm.Unpack,
		cmd_kpm.Inspect,
		cmd_kpm.Run,
		cmd_kpm.New,
		cmd_kpm.Repo,
	},
}
