package flags

import (
	"github.com/rohitramu/kpm/src/cli/model/utils/config"
	"github.com/rohitramu/kpm/src/cli/model/utils/types"
	"github.com/rohitramu/kpm/src/pkg/utils/log"
)

var LogLevel = types.NewFlagBuilder[string]("log-level").
	SetShortDescription("The minimum severity of log messages.").
	SetDefaultValueFunc(func(config *config.KpmConfig) string {
		if result, err := log.DefaultLevel.String(); err != nil {
			log.Panicf("Invalid default log level string: %s", err)
			panic(err)
		} else {
			return result
		}
	}).
	Build()
