package flags

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/constants"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/spf13/pflag"
)

var OutputDirFlag = func() *outputDirFlag {
	var flagName = "output-dir"
	var shortDescription = fmt.Sprintf("Directory in which output files should be written (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<outputName>\" will be deleted if it exists.", constants.GeneratedDirName)

	var workingDir, err = files.GetWorkingDir()
	if err != nil {
		log.Panicf("Failed to get current working directory: %s", err)
	}
	var value = template_package.GetDefaultOutputDir(workingDir)

	var result = &outputDirFlag{
		kpmFlagBase: newKpmFlagBase(
			flagName,
			shortDescription,
			value),
	}

	return result
}()

type outputDirFlag struct {
	*kpmFlagBase[string]
}

func (instance *outputDirFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(
		instance.GetValueRef(),
		instance.GetFlagName(),
		"o",
		instance.GetValue(),
		instance.GetShortDescription())
}
