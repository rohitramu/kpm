package flags

import (
	"github.com/spf13/pflag"
)

var ExportNameFlag = func() *exportNameFlag {
	var flagName = "export-name"
	var shortDescription = "Name of the exported output (defaults to \"<package name>-<package version>\")."
	var value = ""

	var result = &exportNameFlag{
		kpmFlagBase: newKpmFlagBase[string](flagName, shortDescription, value),
	}

	return result
}()

type exportNameFlag struct {
	*kpmFlagBase[string]
}

func (instance *exportNameFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(
		instance.GetValueRef(),
		instance.GetFlagName(),
		"n",
		instance.GetValue(),
		instance.GetShortDescription())
}
