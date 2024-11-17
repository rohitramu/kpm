package cmd_kpm

import (
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

var Version = &types.Command{
	Name:             constants.CmdVersion,
	Alias:            "version",
	ShortDescription: "Prints this KPM binary's version information.",
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		log.Outputf("Version:          %s", constants.VersionString)
		log.Outputf("Git Commit Hash:  %s", constants.GitCommitHash)
		log.Outputf("Build timestamp:  %s", constants.BuildTimestampUTC)

		return nil
	},
}
