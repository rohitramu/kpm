package flags

import (
	"github.com/spf13/pflag"
)

var OutputNameFlag = func() *outputNameFlag {
	var flagName = "output-name"
	var shortDescription = "Name of the output (defaults to \"<package name>-<package version>\")."
	var value = ""

	var result = &outputNameFlag{
		kpmFlagBase: newKpmFlagBase(
			flagName,
			shortDescription,
			value),
	}

	return result
}()

type outputNameFlag struct {
	*kpmFlagBase[string]
}

func (instance *outputNameFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(
		instance.GetValueRef(),
		instance.GetFlagName(),
		"n",
		instance.GetValue(),
		instance.GetShortDescription())
}
