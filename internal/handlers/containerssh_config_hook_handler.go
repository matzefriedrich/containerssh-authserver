package handlers

import (
	"log/slog"

	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/services"

	"github.com/gofiber/fiber/v3"
)

type configHookHandler struct {
	containerAppConfigService services.ContainerAppConfigService
	logger                    *slog.Logger
}

func (h *configHookHandler) Register(app *fiber.App) {
	app.Post("/config", func(ctx fiber.Ctx) error {
		return h.handleUserConfigurationRequest(ctx)
	})
}

var _ RouteHandler = (*configHookHandler)(nil)

// handleUserConfigurationRequest processes user configuration requests and returns Docker container configuration for the authenticated user.
func (h *configHookHandler) handleUserConfigurationRequest(c fiber.Ctx) error {

	c.Accepts("json", "text")

	request := &models.ConfigRequest{}
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	authenticatedUsername := request.AuthenticatedUsername
	cfg, err := h.containerAppConfigService.CreateApplicationConfigFor(authenticatedUsername)
	if err != nil {
		h.logger.Error("Cannot create application config",
			slog.String("username", authenticatedUsername),
			slog.Any("error", err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response := &models.ConfigDockerResponse{
		Username:              request.Username,
		ClientVersion:         request.ClientVersion,
		ConnectionId:          request.ConnectionId,
		RemoteAddress:         request.RemoteAddr,
		AuthenticatedUsername: authenticatedUsername,
		Config:                cfg,
	}

	return c.JSON(response)
}

func NewConfigHookHandler(containerAppConfigService services.ContainerAppConfigService, logger *slog.Logger) RouteHandler {
	return &configHookHandler{
		containerAppConfigService: containerAppConfigService,
		logger:                    logger,
	}
}
