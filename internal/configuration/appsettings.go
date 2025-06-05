package configuration

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type ApplicationConfiguration struct {
	AuthServer AuthServerConfig `mapstructure:"authServer"`
	LogLevel   string           `mapstructure:"logLevel"`
	Port       int              `mapstructure:"port"`
}

type AuthServerConfig struct {
	Users map[string]UserProfile `mapstructure:"users"`
}

type UserProfile struct {
	Binds        []string `mapstructure:"binds"`
	Image        string   `mapstructure:"image"`
	Networks     []string `mapstructure:"networks"`
	PublicKeys   []string `mapstructure:"publicKeys"`
	ShellCommand []string `mapstructure:"shellCommand"`
}

// LoadApplicationSettings populates the provided ApplicationConfiguration instance with values from the app config section.
// Returns an error if the configuration is invalid or cannot be unmarshalled.
func LoadApplicationSettings(validations ...ApplicationConfigurationValidationRule) (*ApplicationConfiguration, error) {

	const applicationSettingsSectionKey = "app"
	options := &ApplicationConfiguration{}
	err := viper.UnmarshalKey(applicationSettingsSectionKey, options)
	if err != nil {
		return nil, errors.New(ErrInvalidApplicationConfiguration)
	}

	validationResultsWithErrors := make([]ValidationResult, 0)
	for _, validation := range validations {
		validationResult, validationErr := validation(options)
		if validationErr != nil {
			return nil, fmt.Errorf("error validating application configuration: %s", validationErr)
		}
		if len(validationResult.errors) > 0 {
			validationResultsWithErrors = append(validationResultsWithErrors, validationResult)
		}
	}
	if len(validationResultsWithErrors) > 0 {
		return nil, NewConfigurationError(ErrInvalidApplicationConfiguration, validationResultsWithErrors)
	}

	return options, nil
}
