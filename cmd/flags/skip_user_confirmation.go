package flags

import (
	"github.com/spf13/pflag"
)

var SkipUserConfirmationFlag = func() *skipUserConfirmationFlag {
	var flagName = "skip-user-confirmation"
	var shortDescription = "Skips user confirmation."
	var value = false

	var result = &skipUserConfirmationFlag{
		kpmFlagBase: newKpmFlagBase[bool](flagName, shortDescription, value),
	}

	return result
}()

type skipUserConfirmationFlag struct {
	*kpmFlagBase[bool]
}

func (instance *skipUserConfirmationFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.BoolVar(
		instance.GetValueRef(),
		instance.GetFlagName(),
		instance.GetValue(),
		instance.GetShortDescription())
}
