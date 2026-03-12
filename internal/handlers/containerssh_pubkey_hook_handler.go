package handlers

import (
	"strings"

	"github.com/rs/zerolog"

	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/services"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/ssh"
)

type pubKeyHookHandler struct {
	profileService services.UserProfileService
	logger         *zerolog.Logger
}

var _ RouteHandler = (*pubKeyHookHandler)(nil)

func (h *pubKeyHookHandler) Register(app *fiber.App) {
	app.Post("/pubkey", h.handlePublicKeyAuthenticationRequest)
}

// handlePublicKeyAuthenticationRequest processes a public key authentication request and verifies the provided public key.
func (h *pubKeyHookHandler) handlePublicKeyAuthenticationRequest(c fiber.Ctx) error {

	c.Accepts("json", "text")

	request := &models.PubKeyRequest{}
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	requestPublicKey, parseErr := parsePublicKey(request.PublicKey)
	if parseErr != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	authorizedKeyString := string(ssh.MarshalAuthorizedKey(requestPublicKey))
	h.logger.Info().
		Str("authorizedKey", authorizedKeyString).
		Msgf("Received public key")

	username := strings.TrimSpace(request.Username)

	publicKeyAccepted, verificationErr := h.profileService.VerifyPublicKey(username, requestPublicKey)
	if verificationErr != nil {
		h.logger.Error().
			Err(verificationErr).
			Str("connectionId", request.ConnectionId).
			Str("username", username).
			Msgf("Public key verification failed")
	}

	if publicKeyAccepted {
		h.logger.Info().
			Str("connectionId", request.ConnectionId).
			Str("username", username).
			Msgf("Public key verification succeeded")
	}

	response := &models.PubKeyResponse{
		AuthenticatedUsername: username,
		ClientVersion:         request.ClientVersion,
		ConnectionId:          request.ConnectionId,
		RemoteAddress:         request.RemoteAddress,
		Success:               publicKeyAccepted,
		Username:              username,
	}

	return c.JSON(response)
}

func parsePublicKey(formattedPublicKey string) (ssh.PublicKey, error) {
	publicKey, _, _, _, parseErr := ssh.ParseAuthorizedKey([]byte(formattedPublicKey))
	if parseErr != nil {
		return nil, parseErr
	}
	return publicKey, nil
}

func NewPubKeyHookHandler(profileService services.UserProfileService, logger *zerolog.Logger) RouteHandler {
	return &pubKeyHookHandler{
		profileService: profileService,
		logger:         logger,
	}
}
