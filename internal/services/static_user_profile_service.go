package services

import (
	"errors"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
)

type staticUserConfigurationProfileService struct {
	settings *configuration.ApplicationConfiguration
}

// GetProfile retrieves the user profile for the given authenticated username from the configuration.
func (u *staticUserConfigurationProfileService) GetProfile(authenticatedUsername string) (configuration.UserProfile, error) {

	profile, found := u.settings.AuthServer.Users[authenticatedUsername]
	if !found {
		return EmptyUserProfile, errors.New(ErrorUserProfileNotFound)
	}

	return profile, nil
}

// RegisterUser is not supported by the static configuration backend as it is read-only.
func (u *staticUserConfigurationProfileService) RegisterUser(_ string, _ configuration.UserProfile) error {
	return errors.New(ErrorReadOnlyBackend)
}

var _ UserProfileBackend = (*staticUserConfigurationProfileService)(nil)

// NewStaticUserConfigurationProfileService creates a new UserProfileBackend using static application configuration.
func NewStaticUserConfigurationProfileService(settings *configuration.ApplicationConfiguration) UserProfileBackend {
	return &staticUserConfigurationProfileService{
		settings: settings,
	}
}
