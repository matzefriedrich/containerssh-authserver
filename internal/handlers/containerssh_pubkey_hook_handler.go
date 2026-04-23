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
	logEvents      *authHookEvents
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

	h.logEvents.PublicKeyReceived(request.ConnectionId, authorizedKeyString)

	username := strings.TrimSpace(request.Username)

	publicKeyAccepted, verificationErr := h.profileService.VerifyPublicKey(username, requestPublicKey)
	if verificationErr != nil {
		h.logEvents.PublicKeyAuthenticationFailed(request.ConnectionId, username, verificationErr)
	}

	h.logEvents.PublicKeyAuthenticationCompeted(request.ConnectionId, username, publicKeyAccepted)

	response := &models.AuthResponse{
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
		logEvents: &authHookEvents{
			logger: logger,
		},
	}
}
