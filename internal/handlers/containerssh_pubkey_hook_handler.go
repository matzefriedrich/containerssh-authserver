package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"golang.org/x/crypto/ssh"
)

type pubKeyHookHandler struct {
	profileService types.UserProfileService
	jsonCodec      JSONCodec
	logEvents      *authHookEvents
}

var _ RouteHandler = (*pubKeyHookHandler)(nil)

func (h *pubKeyHookHandler) Register(router chi.Router) {
	router.Post("/pubkey", h.handlePublicKeyAuthenticationRequest)
}

// handlePublicKeyAuthenticationRequest processes a public key authentication request and verifies the provided public key.
func (h *pubKeyHookHandler) handlePublicKeyAuthenticationRequest(w http.ResponseWriter, r *http.Request) {

	request := &models.PubKeyRequest{}
	if err := h.jsonCodec.Decode(r.Body, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestPublicKey, parseErr := parsePublicKey(request.PublicKey)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorizedKeyString := string(ssh.MarshalAuthorizedKey(requestPublicKey))

	h.logEvents.PublicKeyReceived(request.ConnectionId, authorizedKeyString)

	username := strings.TrimSpace(request.Username)

	ctx := r.Context()
	publicKeyAccepted, verificationErr := h.profileService.VerifyPublicKey(ctx, username, requestPublicKey)
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

	w.Header().Set("Content-Type", "application/json")
	_ = h.jsonCodec.Encode(w, response)
}

func parsePublicKey(formattedPublicKey string) (ssh.PublicKey, error) {
	publicKey, _, _, _, parseErr := ssh.ParseAuthorizedKey([]byte(formattedPublicKey))
	if parseErr != nil {
		return nil, parseErr
	}
	return publicKey, nil
}

func NewPubKeyHookHandler(profileService types.UserProfileService, jsonCodec JSONCodec, logger *slog.Logger) RouteHandler {
	return &pubKeyHookHandler{
		profileService: profileService,
		jsonCodec:      jsonCodec,
		logEvents: &authHookEvents{
			logger: logger,
		},
	}
}
