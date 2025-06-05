package handlers

import (
	"github.com/rs/zerolog"
	"strings"

	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/services"

	"github.com/gofiber/fiber/v2"
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
func (h *pubKeyHookHandler) handlePublicKeyAuthenticationRequest(c *fiber.Ctx) error {

	c.Accepts("json", "text")

	request := &models.PubKeyRequest{}
	if err := c.BodyParser(request); err != nil {
		return err
	}

	requestPublicKey, parseErr := parsePublicKey(request.PublicKey)
	if parseErr != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	h.logger.Info().Msgf("Received public key: %s", string(ssh.MarshalAuthorizedKey(requestPublicKey)))

	username := strings.TrimSpace(request.Username)

	publicKeyAccepted, verificationErr := h.profileService.VerifyPublicKey(username, requestPublicKey)
	if verificationErr != nil {
		h.logger.Error().Msgf("Public key verification failed for user %s: %v", username, verificationErr)
		return c.SendStatus(fiber.StatusForbidden)
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
