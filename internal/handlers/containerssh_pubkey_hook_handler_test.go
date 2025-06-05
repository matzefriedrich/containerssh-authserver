package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers/models"
	"github.com/matzefriedrich/containerssh-authserver/internal/services"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_PubKeyHookHandler_handle_pubkey_request_indicates_success_for_valid_public_key(t *testing.T) {

	// Arrange
	pubKey, _, _ := ed25519.GenerateKey(rand.Reader)
	sshPublicKey, _ := ssh.NewPublicKey(pubKey)
	authorizedKey := string(ssh.MarshalAuthorizedKey(sshPublicKey))
	const username = "john-doe"

	userProfileServiceMock := services.NewUserProfileServiceMock()
	userProfileServiceMock.VerifyPublicKeyFunc = func(authenticatedUsername string, key ssh.PublicKey) (bool, error) {
		keyString := ssh.MarshalAuthorizedKey(key)
		return string(keyString) == authorizedKey &&
			strings.Compare(username, authenticatedUsername) == 0, nil
	}

	sut := NewPubKeyHookHandler(userProfileServiceMock, &zerolog.Logger{})

	requestBody := &models.PubKeyRequest{
		PublicKey: authorizedKey,
		Username:  username,
	}

	body, _ := json.Marshal(requestBody)
	bodyReader := strings.NewReader(string(body))

	app := fiber.New()
	request, _ := http.NewRequest(fiber.MethodPost, "/pubkey", bodyReader)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	request.Header.Set(fiber.HeaderContentLength, fmt.Sprintf("%d", len(body)))

	sut.Register(app)

	// Act
	response, _ := app.Test(request)
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(response.Body)

	responseBody := &models.PubKeyResponse{}
	decoder := json.NewDecoder(response.Body)
	_ = decoder.Decode(responseBody)

	// Assert
	userProfileServiceMock.Verify(services.FunctionVerifyPublicKey, features.TimesOnce(), features.Exact(username), features.Exact(authorizedKey))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, username, responseBody.AuthenticatedUsername)
	assert.Equal(t, username, responseBody.Username)
	assert.True(t, responseBody.Success)
}
