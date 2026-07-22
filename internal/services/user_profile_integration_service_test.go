package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

// Fixed mock
type fixedMockBackend struct {
	getProfile func(string) (configuration.UserProfile, error)
}

func (m *fixedMockBackend) GetProfile(u string) (configuration.UserProfile, error) {
	return m.getProfile(u)
}

func Test_UserProfileIntegrationService_GetProfile_iterates_backends(t *testing.T) {
	// Arrange
	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Backends: []string{"first", "second"},
		},
	}

	backend1 := &fixedMockBackend{
		getProfile: func(string) (configuration.UserProfile, error) {
			return EmptyUserProfile, errors.New(ErrorUserProfileNotFound)
		},
	}
	backend2 := &fixedMockBackend{
		getProfile: func(string) (configuration.UserProfile, error) {
			return configuration.UserProfile{Image: "img2"}, nil
		},
	}

	factory := func(key string) (UserProfileBackend, error) {
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
	profile, err := sut.GetProfile("any")

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
		getProfile: func(u string) (configuration.UserProfile, error) {
			if u == username {
				return configuration.UserProfile{Secret: string(hash), Image: "img"}, nil
			}
			return EmptyUserProfile, errors.New(ErrorUserProfileNotFound)
		},
	}

	factory := func(key string) (UserProfileBackend, error) {
		return backend, nil
	}

	sut := NewUserProfileIntegrationService(factory, settings)

	// Act
	success, err := sut.VerifySecret(username, passwordBase64)

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
		getProfile: func(u string) (configuration.UserProfile, error) {
			if u == username {
				return configuration.UserProfile{PublicKeys: []string{authorizedKey}, Image: "img"}, nil
			}
			return EmptyUserProfile, errors.New(ErrorUserProfileNotFound)
		},
	}

	factory := func(key string) (UserProfileBackend, error) {
		return backend, nil
	}

	sut := NewUserProfileIntegrationService(factory, settings)

	// Act
	success, err := sut.VerifyPublicKey(username, sshPublicKey)

	// Assert
	assert.NoError(t, err)
	assert.True(t, success)
}
