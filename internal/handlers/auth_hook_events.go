package handlers

import "github.com/rs/zerolog"

type authHookEvents struct {
	logger *zerolog.Logger
}

func (e *authHookEvents) PublicKeyAuthenticationCompeted(
	connectionId string,
	username string,
	publicKeyAccepted bool) {
	if !publicKeyAccepted {
		return
	}
	e.logger.Info().
		Str("connectionId", connectionId).
		Str("username", username).
		Msgf("Public key verification succeeded")
}

func (e *authHookEvents) PasswordAuthenticationCompeted(connectionId string, username string, passwordAccepted bool) {
	if !passwordAccepted {
		return
	}
	e.logger.Info().
		Str("connectionId", connectionId).
		Str("username", username).
		Msgf("Password verification succeeded")
}

func (e *authHookEvents) PublicKeyAuthenticationFailed(connectionId string, username string, verificationErr error) {
	e.logger.Error().
		Err(verificationErr).
		Str("connectionId", connectionId).
		Str("username", username).
		Msgf("Public key verification failed")
}

func (e *authHookEvents) PasswordAuthenticationFailed(connectionId string, username string, verificationErr error) {
	e.logger.Error().
		Err(verificationErr).
		Str("connectionId", connectionId).
		Str("username", username).
		Msgf("Password verification failed")

}

func (e *authHookEvents) PublicKeyReceived(connectionId string, authorizedKeyString string) {
	e.logger.Info().
		Str("authorizedKey", authorizedKeyString).
		Msgf("Received public key")

}
