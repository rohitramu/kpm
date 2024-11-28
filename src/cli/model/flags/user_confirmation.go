package flags

import (
	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
)

var UserConfirmation = types.NewFlagBuilder[bool]("confirm").
	SetAlias('y').
	SetShortDescription("Skips user confirmation.").
	SetDefaultValueFunc(func(kc *config.KpmConfig) bool { return false }).
	Build()
