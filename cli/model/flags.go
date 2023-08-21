package model

import (
	"fmt"

	"github.com/rohitramu/kpm/pkg/utils/constants"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
)

var logLevelFlag = NewFlagBuilder[string]("log-level").
	SetShortDescription("The minimum severity of log messages.").
	SetDefaultValue(func() string {
		if result, err := log.DefaultLevel.String(); err != nil {
			log.Panicf("Invalid default log level string: %s", err)
			panic(err)
		} else {
			return result
		}
	}()).
	Build()

var skipUserConfirmationFlag = NewFlagBuilder[bool]("confirm").
	SetShortDescription("Skips user confirmation.").
	SetDefaultValue(false).
	Build()

var packageVersionFlag = NewFlagBuilder[string]("version").
	SetAlias('v').
	SetShortDescription("The template package's version.").
	Build()

var parametersFileFlag = NewFlagBuilder[string]("parameters-file").
	SetAlias('p').
	SetShortDescription("Filepath of the parameters file to use.").
	Build()

var exportDirFlag = NewFlagBuilder[string]("export-dir").
	SetAlias('d').
	SetShortDescription(fmt.Sprintf(
		"The directory which the template package should be exported to (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<export-name>\" will be deleted if it exists.",
		constants.ExportDirName,
	)).
	SetDefaultValue(constants.ExportDirName).
	SetValidationFunc(ValidateDirExists()).
	Build()

var exportNameFlag = NewFlagBuilder[string]("export-name").
	SetAlias('n').
	SetShortDescription("Name of the exported template package (defaults to \"<package name>-<package version>\").").
	Build()

var outputDirFlag = NewFlagBuilder[string]("output-dir").
	SetAlias('d').
	SetShortDescription(fmt.Sprintf(
		"Directory in which output files should be written (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<output-name>\" will be deleted if it exists.",
		constants.GeneratedDirName,
	)).
	SetDefaultValue(func() string {
		var outputDir, err = files.GetAbsolutePath(constants.GeneratedDirName)
		if err != nil {
			log.Panicf("Failed to get default output directory.")
		}

		return outputDir
	}()).
	Build()

var outputNameFlag = NewFlagBuilder[string]("output-name").
	SetAlias('n').
	SetShortDescription("Name of the output (defaults to \"<package name>-<package version>\").").
	Build()

var newPackageOutputDirFlag = NewFlagBuilder[string]("output-dir").
	SetAlias('d').
	SetShortDescription(fmt.Sprintf(
		"Directory in which the new template package should be generated (defaults to \"%s\" under the current working directory) - WARNING: the sub-directory specified by \"<output-name>\" will be deleted if it exists.",
		constants.NewTemplatePackageDirName,
	)).
	SetDefaultValue(func() string {
		var outputDir, err = files.GetAbsolutePath(constants.NewTemplatePackageDirName)
		if err != nil {
			log.Panicf("Failed to get default output directory.")
		}

		return outputDir
	}()).
	Build()
