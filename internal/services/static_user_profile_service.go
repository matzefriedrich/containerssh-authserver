package services

import (
	"context"
	"errors"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
)

type staticUserConfigurationProfileService struct {
	settings *configuration.ApplicationConfiguration
}

// GetProfile retrieves the user profile for the given authenticated username from the configuration.
func (u *staticUserConfigurationProfileService) GetProfile(
	_ context.Context,
	authenticatedUsername string) (types.UserProfile, error) {

	profile, found := u.settings.AuthServer.Users[authenticatedUsername]
	if !found {
		return types.EmptyUserProfile, errors.New(types.ErrorUserProfileNotFound)
	}

	return profile, nil
}

// RegisterUser is not supported by the static configuration backend as it is read-only.
func (u *staticUserConfigurationProfileService) RegisterUser(_ context.Context, _ string, _ types.UserProfile) error {
	return errors.New(types.ErrorReadOnlyBackend)
}

var _ types.UserProfileBackend = (*staticUserConfigurationProfileService)(nil)

// NewStaticUserConfigurationProfileService creates a new UserProfileBackend using static application configuration.
func NewStaticUserConfigurationProfileService(settings *configuration.ApplicationConfiguration) types.UserProfileBackend {
	return &staticUserConfigurationProfileService{
		settings: settings,
	}
}
