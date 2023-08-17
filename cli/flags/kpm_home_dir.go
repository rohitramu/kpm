package flags

import (
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/spf13/pflag"
)

var KpmHomeDirFlag = func() *kpmHomeDirFlag {
	var flagName = "kpm-home-dir"
	var shortDescription = "Directory to use as the KPM home folder."

	value, err := template_package.GetDefaultKpmHomeDir()
	if err != nil {
		log.Panicf("Failed to get default KPM home directory: %s", err)
	}

	var result = &kpmHomeDirFlag{
		kpmFlagBase: newKpmFlagBase[string](flagName, shortDescription, value),
	}

	return result
}()

type kpmHomeDirFlag struct {
	*kpmFlagBase[string]
}

func (instance *kpmHomeDirFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.StringVar(
		instance.GetValueRef(),
		instance.GetFlagName(),
		instance.GetValue(),
		instance.GetShortDescription())
}
