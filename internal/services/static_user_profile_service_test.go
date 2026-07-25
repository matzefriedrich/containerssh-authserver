package services

import (
	"testing"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_StaticUserProfileBackend_GetProfile_returns_correct_profile(t *testing.T) {

	// Arrange
	const username = "johndoe"
	expectedProfile := types.UserProfile{
		PublicKeys: []string{"ssh-rsa ..."},
		Image:      "some-docker-image:tag",
		Secret:     "hash",
	}

	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Users: map[string]types.UserProfile{
				username: expectedProfile,
			},
		},
	}

	sut := NewStaticUserConfigurationProfileService(settings)

	// Act
	actual, err := sut.GetProfile(t.Context(), username)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile, actual)
}

func Test_StaticUserProfileBackend_GetProfile_returns_error_for_unknown_user(t *testing.T) {

	// Arrange
	settings := &configuration.ApplicationConfiguration{
		AuthServer: configuration.AuthServerConfig{
			Users: map[string]types.UserProfile{},
		},
	}

	sut := NewStaticUserConfigurationProfileService(settings)

	// Act
	_, err := sut.GetProfile(t.Context(), "unknown")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, types.ErrorUserProfileNotFound, err.Error())
}

func Test_StaticUserProfileBackend_RegisterUser_returns_readonly_error(t *testing.T) {

	// Arrange
	settings := &configuration.ApplicationConfiguration{}
	sut := NewStaticUserConfigurationProfileService(settings)

	// Act
	err := sut.RegisterUser(t.Context(), "any", types.UserProfile{})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, types.ErrorReadOnlyBackend, err.Error())
}
