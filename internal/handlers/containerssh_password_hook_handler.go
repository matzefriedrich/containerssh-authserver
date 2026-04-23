package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/services"
	"github.com/rs/zerolog"
)

type passwordHookHandler struct {
	profileService services.UserProfileService
	logEvents      *authHookEvents
}

var _ RouteHandler = (*passwordHookHandler)(nil)

func (h *passwordHookHandler) Register(app *fiber.App) {
	app.Post("/password", h.handlePasswordAuthenticationRequest)
}

func (h *passwordHookHandler) handlePasswordAuthenticationRequest(c fiber.Ctx) error {

	c.Accepts("json", "text")

	request := &models.PasswordRequest{}
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	username := strings.TrimSpace(request.Username)

	passwordAccepted, verificationErr := h.profileService.VerifySecret(username, request.PasswordBase64)
	if verificationErr != nil {
		h.logEvents.PasswordAuthenticationFailed(request.ConnectionId, username, verificationErr)
	}

	h.logEvents.PasswordAuthenticationCompeted(request.ConnectionId, username, passwordAccepted)

	response := &models.AuthResponse{
		AuthenticatedUsername: username,
		ClientVersion:         request.ClientVersion,
		ConnectionId:          request.ConnectionId,
		RemoteAddress:         request.RemoteAddress,
		Success:               passwordAccepted,
		Username:              username,
	}

	return c.JSON(response)
}

func NewPasswordHookHandler(profileService services.UserProfileService, logger *zerolog.Logger) RouteHandler {
	return &passwordHookHandler{
		profileService: profileService,
		logEvents: &authHookEvents{
			logger: logger,
		},
	}
}
