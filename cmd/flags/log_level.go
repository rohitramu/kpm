package flags

import (
	"fmt"
	"strings"

	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/spf13/pflag"
)

var LogLevelFlag = func() *logLevelFlag {
	var flagName = "log-level"

	var logNames = make([]string, len(log.LevelNames))
	for i, j := 0, log.LevelNone; i < len(logNames); i, j = i+1, j+1 {
		logNames[i] = log.LevelNames[j]
	}
	var shortDescription = fmt.Sprintf("The minimum severity log level to output - "+
		"the log levels in increasing order of verbosity are: \"%s\"",
		strings.Join(logNames, "\", \""))

	var value = log.DefaultLevel

	var mapping = make(map[log.Level][]string)
	for level, name := range log.LevelNames {
		mapping[level] = []string{name}
	}
	var result = &logLevelFlag{
		kpmEnumFlagBase: newKpmEnumFlagBase[log.Level](flagName, shortDescription, value, mapping),
	}

	return result
}()

type logLevelFlag struct {
	*kpmEnumFlagBase[log.Level]
}

func (instance *logLevelFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.Var(
		instance.GetPFlagValue(),
		instance.GetFlagName(),
		instance.GetShortDescription())
}
