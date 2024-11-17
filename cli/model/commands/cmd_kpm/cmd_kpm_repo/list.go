package cmd_kpm_repo

import (
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/constants"
	"github.com/rohitramu/kpm/cli/model/utils/types"
	"github.com/rohitramu/kpm/pkg"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

var ListCmd = &types.Command{
	Name:             constants.CmdRepoList,
	Alias:            "ls",
	ShortDescription: "Lists the names of available repositories.",
	ExecuteFunc: func(config *config.KpmConfig, args types.ArgCollection) (err error) {
		var repos []string
		repos, err = pkg.ListPackageRepositories(config.Repositories)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			log.Outputf(repo)
		}

		return nil
	},
}
