package flags

import (
	"github.com/spf13/pflag"
)

var PackageVersionFlag = func() *packageVersionFlag {
	var flagName = "package-version"
	var shortDescription = "Version of the package."
	var value = "0.0.0"

	var result = &packageVersionFlag{
		kpmFlagBase: newKpmFlagBase[string](flagName, value, shortDescription),
	}

	return result
}()

type packageVersionFlag struct {
	*kpmFlagBase[string]
}

func (instance *packageVersionFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(
		instance.GetValueRef(),
		instance.GetFlagName(),
		"v",
		instance.GetValue(),
		instance.GetShortDescription())
}
