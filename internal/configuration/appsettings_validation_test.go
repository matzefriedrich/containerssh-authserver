package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EnsurePostgresConfigIfUsed_returns_error_if_postgres_backend_configured_but_postgres_config_is_nil(t *testing.T) {
	// Arrange
	config := &ApplicationConfiguration{
		AuthServer: AuthServerConfig{
			Backends: []string{"postgres"},
		},
		Databases: DatabaseConfig{
			Postgres: nil,
		},
	}
	sut := EnsurePostgresConfigIfUsed()

	// Act
	result, err := sut(config)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrPostgresConfigMissing, err.Error())
	assert.Empty(t, result.errors)
}

func Test_EnsurePostgresConfigIfUsed_returns_no_error_if_postgres_backend_not_configured(t *testing.T) {
	// Arrange
	config := &ApplicationConfiguration{
		AuthServer: AuthServerConfig{
			Backends: []string{"static"},
		},
		Databases: DatabaseConfig{
			Postgres: nil,
		},
	}
	sut := EnsurePostgresConfigIfUsed()

	// Act
	result, err := sut(config)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result.errors)
}

func Test_EnsurePostgresConfigIfUsed_returns_no_error_if_postgres_backend_configured_and_postgres_config_is_valid(t *testing.T) {
	// Arrange
	config := &ApplicationConfiguration{
		AuthServer: AuthServerConfig{
			Backends: []string{"postgres"},
		},
		Databases: DatabaseConfig{
			Postgres: &PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Database: "db",
			},
		},
	}
	sut := EnsurePostgresConfigIfUsed()

	// Act
	result, err := sut(config)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result.errors)
}

func Test_EnsurePostgresConfigIfUsed_returns_validation_errors_if_postgres_config_is_incomplete(t *testing.T) {
	// Arrange
	config := &ApplicationConfiguration{
		Databases: DatabaseConfig{
			Postgres: &PostgresConfig{},
		},
	}
	sut := EnsurePostgresConfigIfUsed()

	// Act
	result, err := sut(config)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.errors, 4)
}
