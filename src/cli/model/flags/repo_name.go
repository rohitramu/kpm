package flags

import (
	"fmt"

	"github.com/rohitramu/kpm/src/cli/model/utils/constants"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
)

var RepoName = types.NewFlagBuilder[string]("repo").
	SetShortDescription(fmt.Sprintf(
		"The repository to interact with.  Defaults to the first repository in the list of available repositories.  Run \"%s %s\" to get the repositories list.",
		constants.CmdRepo,
		constants.CmdRepoList,
	)).
	Build()
