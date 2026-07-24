package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

func Test_PubKeyHookHandler_handle_pubkey_request_indicates_success_for_valid_public_key(t *testing.T) {

	// Arrange
	pubKey, _, _ := ed25519.GenerateKey(rand.Reader)
	sshPublicKey, _ := ssh.NewPublicKey(pubKey)
	authorizedKey := string(ssh.MarshalAuthorizedKey(sshPublicKey))
	const username = "john-doe"

	userProfileServiceMock := types.NewUserProfileServiceMock()
	userProfileServiceMock.VerifyPublicKeyFunc = func(_ context.Context, authenticatedUsername string, key ssh.PublicKey) (bool, error) {
		keyString := ssh.MarshalAuthorizedKey(key)
		return string(keyString) == authorizedKey &&
			strings.Compare(username, authenticatedUsername) == 0, nil
	}

	handler := slog.NewTextHandler(io.Discard, nil)
	logger := slog.New(handler)
	jsonCodec := NewJSONCodec()
	sut := NewPubKeyHookHandler(userProfileServiceMock, jsonCodec, logger)

	requestBody := &models.PubKeyRequest{
		PublicKey: authorizedKey,
		Username:  username,
	}

	body, _ := json.Marshal(requestBody)
	bodyReader := strings.NewReader(string(body))

	router := chi.NewRouter()
	sut.Register(router)

	request, _ := http.NewRequest(http.MethodPost, "/pubkey", bodyReader)
	request.Header.Set("Content-Type", "application/json")

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(response.Body)

	responseBody := &models.AuthResponse{}
	decoder := json.NewDecoder(response.Body)
	_ = decoder.Decode(responseBody)

	// Assert
	userProfileServiceMock.Verify(types.Function_UserProfileService_VerifyPublicKey, features.TimesOnce(), features.Exact(username), features.Exact(authorizedKey))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, username, responseBody.AuthenticatedUsername)
	assert.Equal(t, username, responseBody.Username)
	assert.True(t, responseBody.Success)
}
