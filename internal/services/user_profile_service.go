package services

import (
	"bytes"
	"errors"
	"strings"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"golang.org/x/crypto/ssh"
)

const (
	ErrorUserProfileNotFound          = "user profile not found"
	ErrorInvalidSessionContainerImage = "invalid session container image"
)

var (
	EmptyUserProfile = configuration.UserProfile{}
)

type staticUserConfigurationProfileService struct {
	settings *configuration.ApplicationConfiguration
}

// VerifyPublicKey compares the given public key with stored keys for the specified user and returns true if a match is found.
func (u *staticUserConfigurationProfileService) VerifyPublicKey(username string, expectedKey ssh.PublicKey) (bool, error) {

	profile, err := u.GetProfile(username)
	if err != nil {
		return false, err
	}

	expectedKeyBytes := expectedKey.Marshal()
	for _, formattedPublicKey := range profile.PublicKeys {
		key, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(formattedPublicKey))
		if bytes.Equal(key.Marshal(), expectedKeyBytes) {
			return true, nil
		}
	}

	return false, nil
}

// GetProfile retrieves the user profile for the given authenticated username from the configuration.
func (u *staticUserConfigurationProfileService) GetProfile(authenticatedUsername string) (configuration.UserProfile, error) {

	profile, found := u.settings.AuthServer.Users[authenticatedUsername]
	if found == false {
		return EmptyUserProfile, errors.New(ErrorUserProfileNotFound)
	}

	containerImageName := strings.TrimSpace(profile.Image)
	if containerImageName == "" {
		return EmptyUserProfile, errors.New(ErrorInvalidSessionContainerImage)
	}

	shellCommand := profile.ShellCommand
	if len(shellCommand) == 0 {
		profile.ShellCommand = []string{"/bin/sh"}
	}

	return profile, nil
}

var _ UserProfileService = (*staticUserConfigurationProfileService)(nil)

// NewStaticUserConfigurationProfileService creates a new UserProfileService using static application configuration.
func NewStaticUserConfigurationProfileService(settings *configuration.ApplicationConfiguration) UserProfileService {
	return &staticUserConfigurationProfileService{
		settings: settings,
	}
}
