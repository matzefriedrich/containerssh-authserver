package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
)

type passwordHookHandler struct {
	profileService types.UserProfileService
	jsonCodec      JSONCodec
	logEvents      *authHookEvents
}

var _ RouteHandler = (*passwordHookHandler)(nil)

func (h *passwordHookHandler) Register(router chi.Router) {
	router.Post("/password", h.handlePasswordAuthenticationRequest)
}

func (h *passwordHookHandler) handlePasswordAuthenticationRequest(w http.ResponseWriter, r *http.Request) {

	request := &models.PasswordRequest{}
	if err := h.jsonCodec.Decode(r.Body, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(request.Username)

	ctx := r.Context()
	passwordAccepted, verificationErr := h.profileService.VerifySecret(ctx, username, request.PasswordBase64)
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

	w.Header().Set("Content-Type", "application/json")
	_ = h.jsonCodec.Encode(w, response)
	w.WriteHeader(http.StatusOK)
}

func NewPasswordHookHandler(profileService types.UserProfileService, jsonCodec JSONCodec, logger *slog.Logger) RouteHandler {
	return &passwordHookHandler{
		profileService: profileService,
		jsonCodec:      jsonCodec,
		logEvents: &authHookEvents{
			logger: logger,
		},
	}
}
