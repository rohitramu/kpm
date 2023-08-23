package model

import (
	"fmt"
	"path/filepath"

	"github.com/caarlos0/env/v9"
	"gopkg.in/yaml.v3"

	"github.com/rohitramu/kpm/pkg/utils/constants"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
	"github.com/rohitramu/kpm/pkg/utils/template_repository"
)

// TODO: Reduce amount of error wrapping.

type KpmConfig struct {
	LogLevel     log.Level                                `yaml:"logLevel"     env:"LOG_LEVEL"`
	Repositories template_repository.RepositoryCollection `yaml:"repositories" env:"REPOSITORIES"`
}

func ReadConfig() (*KpmConfig, error) {
	var err error

	// Set defaults
	var result = &KpmConfig{
		LogLevel:     log.DefaultLevel,
		Repositories: template_repository.RepositoryCollection{},
	}

	var kpmHomeDir string
	kpmHomeDir, err = template_package.GetKpmHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to locate KPM home directory: %s", err)
	}

	// Get values from KPM home directory's config file.
	err = result.readConfigFromFile(filepath.Join(kpmHomeDir, constants.KpmConfigFileName))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read KPM configuration file in the KPM home directory '%s': %s",
			kpmHomeDir,
			err,
		)
	}

	// Merge in values from current working directory's config file.
	var workingDir string
	workingDir, err = files.GetWorkingDir()
	if err != nil {
		log.Panicf("Failed to get current working directory: %s", err)
	}
	err = result.readConfigFromFile(filepath.Join(workingDir, constants.KpmConfigFileName))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read KPM configuration file in the current working directory '%s': %s",
			workingDir,
			err,
		)
	}

	// Merge in values from environment variables.
	err = env.ParseWithOptions(result, env.Options{Prefix: constants.KpmEnvVariablePrefix})
	if err != nil {
		return nil, fmt.Errorf("failed to parse KPM configuration from environment variables: %s", err)
	}

	return result, nil
}

func (result *KpmConfig) readConfigFromFile(filepath string) (err error) {
	// Get the absolute filepath.
	var absoluteFilePath string
	absoluteFilePath, err = files.GetAbsolutePath(filepath)
	if err != nil {
		return fmt.Errorf("invalid path '%s' to config file: %s", filepath, err)
	}

	// Make sure the file exists.
	err = files.FileExists(absoluteFilePath, "KPM configuration")
	if err != nil {
		log.Verbosef("KPM config file at '%s' doesn't exist: %s", absoluteFilePath, err)
		return nil
	}

	// Read file.
	var yamlBytes []byte
	yamlBytes, err = files.ReadBytes(absoluteFilePath)
	if err != nil {
		return fmt.Errorf("failed to read KPM configuration file at '%s': %s", absoluteFilePath, err)
	}

	// Unmarshal the YAML into an object.
	// TODO: Use my own yaml unmarshaller once it works with strict unmarshalling: https://github.com/go-yaml/yaml/issues/460
	err = yaml.Unmarshal(yamlBytes, result)
	if err != nil {
		return fmt.Errorf("failed to parse configuration file at '%s': %s", absoluteFilePath, err)
	}

	return nil
}
