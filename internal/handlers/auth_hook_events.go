package handlers

import "log/slog"

type authHookEvents struct {
	logger *slog.Logger
}

func (e *authHookEvents) PublicKeyAuthenticationCompeted(
	connectionId string,
	username string,
	publicKeyAccepted bool) {
	if !publicKeyAccepted {
		return
	}
	e.logger.Info("Public key verification succeeded",
		slog.String("connectionId", connectionId),
		slog.String("username", username))
}

func (e *authHookEvents) PasswordAuthenticationCompeted(connectionId string, username string, passwordAccepted bool) {
	if !passwordAccepted {
		return
	}
	e.logger.Info("Password verification succeeded",
		slog.String("connectionId", connectionId),
		slog.String("username", username))
}

func (e *authHookEvents) PublicKeyAuthenticationFailed(connectionId string, username string, verificationErr error) {
	e.logger.Error("Public key verification failed",
		"error", verificationErr,
		slog.String("connectionId", connectionId),
		slog.String("username", username))
}

func (e *authHookEvents) PasswordAuthenticationFailed(connectionId string, username string, verificationErr error) {
	e.logger.Error("Password verification failed",
		"error", verificationErr,
		slog.String("connectionId", connectionId),
		slog.String("username", username))
}

func (e *authHookEvents) PublicKeyReceived(connectionId string, authorizedKeyString string) {
	e.logger.Info("Received public key",
		slog.String("connectionId", connectionId),
		slog.String("authorizedKey", authorizedKeyString))
}
