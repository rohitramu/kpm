package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/rohitramu/kpm/cli/commands"
	"github.com/rohitramu/kpm/pkg/common"
	"github.com/rohitramu/kpm/pkg/utils/files"
	"github.com/rohitramu/kpm/pkg/utils/log"
	"github.com/rohitramu/kpm/pkg/utils/template_package"
)

const ConfigFileName = ".kpm-config"
const ConfigFileExt = "yaml"
const KpmHomeDirEnvVar = "KPM_HOME"

func InitConfig(conf *common.KpmConfigSchema) error {
	var err error

	err = viper.BindPFlags(commands.KpmCmd.PersistentFlags())
	if err != nil {
		log.Panicf("failed to bind config to flags: %s", err)
	}

	var kpmHomeDir string
	kpmHomeDir, err = getKpmHomeDirFromEnv()
	if err != nil {
		return fmt.Errorf("failed to get KPM home directory: %s", err)
	}
	var kpmHomeDirAbs string
	kpmHomeDirAbs, err = files.GetAbsolutePath(kpmHomeDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of KPM home directory '%s': %s", kpmHomeDir, err)
	}

	viper.SetConfigType(ConfigFileExt)

	// Read config from KPM home directory
	viper.SetConfigFile(filepath.Join(kpmHomeDirAbs, ConfigFileName))
	err = viper.ReadInConfig()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Panicf("Failed to read config from '%s': %s", kpmHomeDirAbs, err)
	}

	// Read config from current directory
	var pwd string
	pwd, err = files.GetWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %s", err)
	}
	viper.SetConfigFile(filepath.Join(pwd, ConfigFileName))
	err = viper.MergeInConfig()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Panicf("Failed to read config from current working directory '%s': %s", pwd, err)
	}

	var settings = viper.AllSettings()
	log.Debugf("%v", settings)
	err = viper.UnmarshalExact(conf)
	if err != nil {
		return fmt.Errorf("invalid KPM config: %s", err)
	}

	return nil
}

func getKpmHomeDirFromEnv() (string, error) {
	var err error

	// Try to get the KPM home directory from the environment variable.
	var kpmHomeDir = strings.TrimSpace(os.ExpandEnv("$" + KpmHomeDirEnvVar))
	if kpmHomeDir != "" {
		var kpmHomeDirAbs string
		kpmHomeDirAbs, err = files.GetAbsolutePath(kpmHomeDir)
		if err != nil {
			return "", fmt.Errorf("invalid directory specified for the \"KPM_HOME\" environment variable '%s': %s", kpmHomeDir, err)
		}

		kpmHomeDir = kpmHomeDirAbs
	} else {
		// Since the environment variable was empty or not defined, use the default value.
		kpmHomeDir, err = template_package.GetDefaultKpmHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get default KPM home directory: %s", err)
		}
	}

	return kpmHomeDir, nil
}
