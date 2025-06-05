package configuration

import (
	"errors"
	"golang.org/x/crypto/ssh"
)

type ApplicationConfigurationFunc func(config *ApplicationConfiguration) error

// EnsureUnprivilegedPortOrDefault ensures the application's port is above 1024, assigning a default if the configured port is privileged.
func EnsureUnprivilegedPortOrDefault(defaultPort int) ApplicationConfigurationValidationRule {
	return func(config *ApplicationConfiguration) (ValidationResult, error) {
		result := ValidationResult{}
		if config.Port < 1024 {
			if defaultPort < 1024 {
				result.errors = append(result.errors, errors.New(ErrUnprivilegedPortExpected))
			}
			config.Port = defaultPort
		}
		return result, nil
	}
}

// EnsurePublicKeyFormat validates that all user public keys in the application configuration are in a proper authorized key format.
func EnsurePublicKeyFormat() ApplicationConfigurationValidationRule {
	return func(config *ApplicationConfiguration) (ValidationResult, error) {
		result := ValidationResult{}
		if config.AuthServer.Users == nil {
			return result, nil
		}
		for _, profile := range config.AuthServer.Users {
			for _, formattedPublicKey := range profile.PublicKeys {
				_, _, _, _, parseErr := ssh.ParseAuthorizedKey([]byte(formattedPublicKey))
				if parseErr != nil {
					result.errors = append(result.errors, errors.New(ErrInvalidPublicKeyFormat))
				}
			}
		}
		return result, nil
	}
}

// EnsureUserProfilesNotEmpty validates that the AuthServer.Users section in the configuration is not nil or empty.
func EnsureUserProfilesNotEmpty() ApplicationConfigurationValidationRule {
	return func(config *ApplicationConfiguration) (ValidationResult, error) {
		if config.AuthServer.Users == nil || len(config.AuthServer.Users) == 0 {
			return ValidationResult{}, errors.New(ErrNoUsersConfigured)
		}
		return ValidationResult{}, nil
	}
}

type ValidationResult struct {
	errors []error
}

type ApplicationConfigurationValidationRule func(config *ApplicationConfiguration) (ValidationResult, error)
