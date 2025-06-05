package services

import (
	"crypto/rand"
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	"testing"
)

func Test_UserProfileService_VerifyPublicKey_returns_true_for_known_user_and_matching_key(t *testing.T) {

	// Arrange
	pubKey, _, _ := ed25519.GenerateKey(rand.Reader)
	sshPublicKey, _ := ssh.NewPublicKey(pubKey)

	authorizedKey := string(ssh.MarshalAuthorizedKey(sshPublicKey))
	username := "johndoe"

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
