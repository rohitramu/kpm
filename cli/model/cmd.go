package model

import (
	"github.com/rohitramu/kpm/cli/model/utils/constants"
)

var KpmCmd = &Command{
	Name: constants.CmdKpm,
	Flags: FlagCollection{
		StringFlags: []Flag[string]{
			logLevelFlag,
		},
	},
	SubCommands: []*Command{
		listCmd,
		removeCmd,
		purgeCmd,
		packCmd,
		unpackCmd,
		inspectCmd,
		runCmd,
		newCmd,
		repoCmd,
	},
}
