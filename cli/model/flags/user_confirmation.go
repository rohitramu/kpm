package flags

import (
	"github.com/rohitramu/kpm/cli/model/utils/config"
	"github.com/rohitramu/kpm/cli/model/utils/types"
)

var UserConfirmation = types.NewFlagBuilder[bool]("confirm").
	SetShortDescription("Skips user confirmation.").
	SetDefaultValueFunc(func(kc *config.KpmConfig) bool { return false }).
	Build()
