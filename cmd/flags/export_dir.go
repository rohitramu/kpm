package flags

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/constants"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/spf13/pflag"
)

var ExportDirFlag = func() *exportDirFlag {
	var flagName = "export-dir"
	var shortDescription = fmt.Sprintf("Directory in which exported files should be written (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<exportName>\" will be deleted if it exists.", constants.ExportDirName)

	var workingDir string
	var err error
	workingDir, err = files.GetWorkingDir()
	if err != nil {
		log.Panicf("failed to get current working directory: %s", err)
	}
	var value = template_package.GetDefaultExportDir(workingDir)

	var result = &exportDirFlag{
		kpmFlagBase: newKpmFlagBase[string](flagName, shortDescription, value),
	}

	return result
}()

type exportDirFlag struct {
	*kpmFlagBase[string]
}

func (instance *exportDirFlag) Add(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(
		instance.GetValueRef(),
		instance.GetFlagName(),
		"d",
		instance.GetValue(),
		instance.GetShortDescription())
}
