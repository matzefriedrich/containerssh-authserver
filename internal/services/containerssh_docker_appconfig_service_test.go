package services

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/shims"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/stretchr/testify/assert"
)

func Test_CreateApplicationConfigFor_returns_error_if_user_not_found(t *testing.T) {
	// Arrange
	profileServiceMock := NewUserProfileServiceMock()
	profileServiceMock.GetProfileFunc = func(authenticatedUser string) (configuration.UserProfile, error) {
		return EmptyUserProfile, errors.New(ErrorUserProfileNotFound)
	}

	const authenticatedUser = "johndoe"
	sut := NewDockerAppConfigService(profileServiceMock, slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Act
	actual, err := sut.CreateApplicationConfigFor(authenticatedUser)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, shims.AppConfigShim{}, actual)

	profileServiceMock.Verify(Function_UserProfileService_GetProfile, features.TimesOnce(), features.Exact(authenticatedUser))
}

func Test_CreateApplicationConfigFor_returns_docker_container_configuration_for_known_user(t *testing.T) {
	// Arrange
	expectedDockerImage := "some-docker-image:tag"
	expectedContainerNetworkName := "some-network"

	profileServiceMock := NewUserProfileServiceMock()
	profileServiceMock.GetProfileFunc = func(authenticatedUser string) (configuration.UserProfile, error) {
		return configuration.UserProfile{
			Image:        expectedDockerImage,
			ShellCommand: []string{"/bin/sh"},
			Networks:     []string{expectedContainerNetworkName},
			Binds:        []string{"/some/host/path:/some/container/path"},
		}, nil
	}

	const authenticatedUser = "johndoe"
	sut := NewDockerAppConfigService(profileServiceMock, slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Act
	actual, err := sut.CreateApplicationConfigFor(authenticatedUser)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedDockerImage, actual.Docker.Execution.ContainerConfig.Image)

	_, found := actual.Docker.Execution.NetworkConfig.EndpointsConfig[expectedContainerNetworkName]
	assert.True(t, found)

	profileServiceMock.Verify(Function_UserProfileService_GetProfile, features.TimesOnce(), features.Exact(authenticatedUser))
}
