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
	Databases  DatabaseConfig   `mapstructure:"databases"`
}

type AuthServerConfig struct {
	Backends []string               `mapstructure:"backends"`
	Users    map[string]UserProfile `mapstructure:"users"`
}

type DatabaseConfig struct {
	Postgres *PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (p PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode)
}

type UserProfile struct {
	Binds        []string `mapstructure:"binds"`
	Image        string   `mapstructure:"image"`
	Networks     []string `mapstructure:"networks"`
	PublicKeys   []string `mapstructure:"publicKeys"`
	ShellCommand []string `mapstructure:"shellCommand"`
	Secret       string   `mapstructure:"secret"`
}

// LoadApplicationSettings populates the provided ApplicationConfiguration instance with values from the app config section.
// Returns an error if the configuration is invalid or cannot be unmarshaled.
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
