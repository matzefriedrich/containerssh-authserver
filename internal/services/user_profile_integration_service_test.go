package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

// Fixed mock
type fixedMockBackend struct {
	getProfile   func(context.Context, string) (types.UserProfile, error)
	registerUser func(context.Context, string, types.UserProfile) error
}

func (m *fixedMockBackend) GetProfile(ctx context.Context, u string) (types.UserProfile, error) {
	return m.getProfile(ctx, u)
}

func (m *fixedMockBackend) RegisterUser(ctx context.Context, username string, profile types.UserProfile) error {
	if m.registerUser != nil {
		return m.registerUser(ctx, username, profile)
	}
	return nil
}

func Test_UserProfileIntegrationService_GetProfile_iterates_backends(t *testing.T) {
	// Arrange
	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Backends: []string{"first", "second"},
		},
	}

	backend1 := &fixedMockBackend{
		getProfile: func(context.Context, string) (types.UserProfile, error) {
			return types.EmptyUserProfile, errors.New(types.ErrorUserProfileNotFound)
		},
	}
	backend2 := &fixedMockBackend{
		getProfile: func(context.Context, string) (types.UserProfile, error) {
			return types.UserProfile{Image: "img2"}, nil
		},
	}

	factory := func(key string) (types.UserProfileBackend, error) {
		if key == "first" {
			return backend1, nil
		}
		if key == "second" {
			return backend2, nil
		}
		return nil, errors.New("unexpected backend key")
	}

	sut := NewUserProfileIntegrationService(factory, settings)

	// Act
	profile, err := sut.GetProfile(t.Context(), "any")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "img2", profile.Image)
}

func Test_UserProfileIntegrationService_VerifySecret(t *testing.T) {
	// Arrange
	const username = "user1"
	const password = "password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password))

	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Backends: []string{"static"},
		},
	}

	backend := &fixedMockBackend{
		getProfile: func(ctx context.Context, u string) (types.UserProfile, error) {
			if u == username {
				return types.UserProfile{Secret: string(hash), Image: "img"}, nil
			}
			return types.EmptyUserProfile, errors.New(types.ErrorUserProfileNotFound)
		},
	}

	factory := func(key string) (types.UserProfileBackend, error) {
		return backend, nil
	}

	sut := NewUserProfileIntegrationService(factory, settings)

	// Act
	success, err := sut.VerifySecret(t.Context(), username, passwordBase64)

	// Assert
	assert.NoError(t, err)
	assert.True(t, success)
}

func Test_UserProfileIntegrationService_VerifyPublicKey(t *testing.T) {
	// Arrange
	const username = "user1"
	pubKey, _, _ := ed25519.GenerateKey(rand.Reader)
	sshPublicKey, _ := ssh.NewPublicKey(pubKey)
	authorizedKey := string(ssh.MarshalAuthorizedKey(sshPublicKey))

	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Backends: []string{"static"},
		},
	}

	backend := &fixedMockBackend{
		getProfile: func(ctx context.Context, u string) (types.UserProfile, error) {
			if u == username {
				return types.UserProfile{PublicKeys: []string{authorizedKey}, Image: "img"}, nil
			}
			return types.EmptyUserProfile, errors.New(types.ErrorUserProfileNotFound)
		},
	}

	factory := func(key string) (types.UserProfileBackend, error) {
		return backend, nil
	}

	sut := NewUserProfileIntegrationService(factory, settings)

	// Act
	success, err := sut.VerifyPublicKey(t.Context(), username, sshPublicKey)

	// Assert
	assert.NoError(t, err)
	assert.True(t, success)
}

func Test_UserProfileIntegrationService_RegisterUser_skips_readonly_backends(t *testing.T) {
	// Arrange
	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Backends: []string{"readonly", "writable"},
		},
	}

	backend1 := &fixedMockBackend{
		registerUser: func(context.Context, string, types.UserProfile) error {
			return errors.New(types.ErrorReadOnlyBackend)
		},
	}
	backend2 := &fixedMockBackend{
		registerUser: func(context.Context, string, types.UserProfile) error {
			return nil
		},
	}

	factory := func(key string) (types.UserProfileBackend, error) {
		if key == "readonly" {
			return backend1, nil
		}
		if key == "writable" {
			return backend2, nil
		}
		return nil, errors.New("unexpected backend key")
	}

	sut := NewUserProfileIntegrationService(factory, settings)

	// Act
	err := sut.RegisterUser(t.Context(), "user", types.UserProfile{})

	// Assert
	assert.NoError(t, err)
}

func Test_UserProfileIntegrationService_RegisterUser_returns_error_if_no_writable_backend(t *testing.T) {
	// Arrange
	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Backends: []string{"readonly"},
		},
	}

	backend := &fixedMockBackend{
		registerUser: func(context.Context, string, types.UserProfile) error {
			return errors.New(types.ErrorReadOnlyBackend)
		},
	}

	factory := func(key string) (types.UserProfileBackend, error) {
		return backend, nil
	}

	sut := NewUserProfileIntegrationService(factory, settings)

	// Act
	err := sut.RegisterUser(t.Context(), "user", types.UserProfile{})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "no writable backend found", err.Error())
}
