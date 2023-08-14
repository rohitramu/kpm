package flags

import (
	"github.com/spf13/pflag"
)

var ParametersFileFlag = func() *parametersFileFlag {
	var flagName = "parameters-file"
	var shortDescription = "Filepath of the parameters file to use."
	var value = ""

	var result = &parametersFileFlag{
		kpmFlagBase: newKpmFlagBase[string](flagName, shortDescription, value),
	}

	return result
}()

type parametersFileFlag struct {
	*kpmFlagBase[string]
}

func (instance *parametersFileFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(
		instance.GetValueRef(),
		instance.GetFlagName(),
		"f",
		instance.GetValue(),
		instance.GetShortDescription())
}
