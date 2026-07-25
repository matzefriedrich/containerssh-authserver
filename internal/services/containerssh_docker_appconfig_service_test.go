package services

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/matzefriedrich/containerssh-authserver/internal/shims"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/stretchr/testify/assert"
)

func Test_CreateApplicationConfigFor_returns_error_if_user_not_found(t *testing.T) {
	// Arrange
	profileServiceMock := types.NewUserProfileServiceMock()
	profileServiceMock.GetProfileFunc = func(ctx context.Context, authenticatedUser string) (types.UserProfile, error) {
		return types.EmptyUserProfile, errors.New(types.ErrorUserProfileNotFound)
	}

	const authenticatedUser = "johndoe"
	sut := NewDockerAppConfigService(profileServiceMock, slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Act
	actual, err := sut.CreateApplicationConfigFor(t.Context(), authenticatedUser)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, shims.AppConfigShim{}, actual)

	profileServiceMock.Verify(types.Function_UserProfileService_GetProfile, features.TimesOnce(), features.Exact(authenticatedUser))
}

func Test_CreateApplicationConfigFor_returns_docker_container_configuration_for_known_user(t *testing.T) {
	// Arrange
	expectedDockerImage := "some-docker-image:tag"
	expectedContainerNetworkName := "some-network"

	profileServiceMock := types.NewUserProfileServiceMock()
	profileServiceMock.GetProfileFunc = func(ctx context.Context, authenticatedUser string) (types.UserProfile, error) {
		return types.UserProfile{
			Image:        expectedDockerImage,
			ShellCommand: []string{"/bin/sh"},
			Networks:     []string{expectedContainerNetworkName},
			Binds:        []string{"/some/host/path:/some/container/path"},
		}, nil
	}

	const authenticatedUser = "johndoe"
	sut := NewDockerAppConfigService(profileServiceMock, slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Act
	actual, err := sut.CreateApplicationConfigFor(t.Context(), authenticatedUser)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedDockerImage, actual.Docker.Execution.ContainerConfig.Image)

	_, found := actual.Docker.Execution.NetworkConfig.EndpointsConfig[expectedContainerNetworkName]
	assert.True(t, found)

	profileServiceMock.Verify(types.Function_UserProfileService_GetProfile, features.TimesOnce(), features.Exact(authenticatedUser))
}
