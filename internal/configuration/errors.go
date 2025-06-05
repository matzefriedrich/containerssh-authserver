package configuration

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
)

const (
	ErrInvalidApplicationConfiguration = "invalid application configuration"
	ErrUnprivilegedPortExpected        = "port must be greater than 1024"
	ErrInvalidPublicKeyFormat          = "invalid public key format"
	ErrNoUsersConfigured               = "no users configured"
)

type Error struct {
	types.AuthServerError
	Results []ValidationResult
}

// NewConfigurationError creates an error with a message, validation results, and initializes it using the provided functions.
func NewConfigurationError(msg string, validationResults []ValidationResult, initializers ...types.AuthServerErrorFunc) error {
	err := &Error{
		AuthServerError: types.AuthServerError{Msg: msg},
		Results:         validationResults,
	}
	for _, initializer := range initializers {
		initializer(&err.AuthServerError)
		initializer(err)
	}
	return err
}
