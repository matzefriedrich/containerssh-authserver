package configuration

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

const (
	appSettingsConfigName = "config"
)

// PathOption adds a configuration path to viper for loading configuration files.
func PathOption(path string) func() {
	return func() {
		viper.AddConfigPath(path)
	}
}

// ConfigTypesOption sets multiple configuration types in viper.
func ConfigTypesOption(configType ...string) func() {
	return func() {
		for _, configType := range configType {
			viper.SetConfigType(configType)
		}
	}
}

// ConfigureApplication Configures and reads the application configuration.
func ConfigureApplication(options ...func()) error {

	for _, option := range options {
		option()
	}

	dotReplacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(dotReplacer)
	viper.AutomaticEnv()

	configFiles := []string{appSettingsConfigName}

	for _, configFile := range configFiles {
		viper.SetConfigName(configFile)
		err := viper.MergeInConfig()
		if err != nil {
			return fmt.Errorf("error reading config file: %s", configFile)
		}
	}

	return nil
}
