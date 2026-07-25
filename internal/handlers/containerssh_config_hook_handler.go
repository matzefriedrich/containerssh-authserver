package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
)

type configHookHandler struct {
	containerAppConfigService types.ContainerAppConfigService
	jsonCodec                 JSONCodec
	logger                    *slog.Logger
}

func (h *configHookHandler) Register(router chi.Router) {
	router.Post("/config", h.handleUserConfigurationRequest)
}

var _ RouteHandler = (*configHookHandler)(nil)

// handleUserConfigurationRequest processes user configuration requests and returns Docker container configuration for the authenticated user.
func (h *configHookHandler) handleUserConfigurationRequest(w http.ResponseWriter, r *http.Request) {

	request := &models.ConfigRequest{}
	if err := h.jsonCodec.Decode(r.Body, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authenticatedUsername := request.AuthenticatedUsername
	cfg, err := h.containerAppConfigService.CreateApplicationConfigFor(r.Context(), authenticatedUsername)
	if err != nil {
		h.logger.Error("Cannot create application config",
			slog.String("username", authenticatedUsername),
			slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &models.ConfigDockerResponse{
		Username:              request.Username,
		ClientVersion:         request.ClientVersion,
		ConnectionId:          request.ConnectionId,
		RemoteAddress:         request.RemoteAddr,
		AuthenticatedUsername: authenticatedUsername,
		Config:                cfg,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = h.jsonCodec.Encode(w, response)
	w.WriteHeader(http.StatusOK)
}

func NewConfigHookHandler(containerAppConfigService types.ContainerAppConfigService, jsonCodec JSONCodec, logger *slog.Logger) RouteHandler {
	return &configHookHandler{
		containerAppConfigService: containerAppConfigService,
		jsonCodec:                 jsonCodec,
		logger:                    logger,
	}
}
