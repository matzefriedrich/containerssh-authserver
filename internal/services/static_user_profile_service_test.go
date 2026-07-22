package services

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

func Test_UserProfileService_VerifyPublicKey_returns_true_for_known_user_and_matching_key(t *testing.T) {

	// Arrange
	pubKey, _, _ := ed25519.GenerateKey(rand.Reader)
	sshPublicKey, _ := ssh.NewPublicKey(pubKey)

	authorizedKey := string(ssh.MarshalAuthorizedKey(sshPublicKey))
	const username = "johndoe"

	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Users: map[string]configuration.UserProfile{
				username: {
					PublicKeys: []string{authorizedKey},
					Image:      "some-docker-image:tag",
				},
			},
		},
	}

	sut := NewStaticUserConfigurationProfileService(settings)

	// Act
	actual, err := sut.VerifyPublicKey(username, sshPublicKey)

	// Assert
	assert.NoError(t, err)
	assert.True(t, actual)
}

func Test_UserProfileService_VerifySecret_returns_true_for_known_user_and_matching_password(t *testing.T) {

	// Arrange
	const username = "johndoe"
	const password = "secret-password"
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password))

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Users: map[string]configuration.UserProfile{
				username: {
					Secret: string(hash),
					Image:  "some-docker-image:tag",
				},
			},
		},
	}

	sut := NewStaticUserConfigurationProfileService(settings)

	// Act
	actual, err := sut.VerifySecret(username, passwordBase64)

	// Assert
	assert.NoError(t, err)
	assert.True(t, actual)
}

func Test_UserProfileService_VerifySecret_returns_false_for_incorrect_password(t *testing.T) {

	// Arrange
	const username = "johndoe"
	const password = "secret-password"
	const wrongPassword = "wrong-password"
	wrongPasswordBase64 := base64.StdEncoding.EncodeToString([]byte(wrongPassword))

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Users: map[string]configuration.UserProfile{
				username: {
					Secret: string(hash),
					Image:  "some-docker-image:tag",
				},
			},
		},
	}

	sut := NewStaticUserConfigurationProfileService(settings)

	// Act
	actual, err := sut.VerifySecret(username, wrongPasswordBase64)

	// Assert
	assert.Error(t, err)
	assert.False(t, actual)
	assert.Equal(t, "password invalid", err.Error())
}

func Test_UserProfileService_VerifySecret_returns_error_for_unknown_user(t *testing.T) {

	// Arrange
	const username = "johndoe"
	const unknownUser = "unknown"
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte("any"))

	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Users: map[string]configuration.UserProfile{
				username: {
					Secret: "any",
					Image:  "some-docker-image:tag",
				},
			},
		},
	}

	sut := NewStaticUserConfigurationProfileService(settings)

	// Act
	actual, err := sut.VerifySecret(unknownUser, passwordBase64)

	// Assert
	assert.Error(t, err)
	assert.False(t, actual)
	assert.Contains(t, err.Error(), "failed to get user profile")
}
