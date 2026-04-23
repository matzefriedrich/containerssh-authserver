package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"golang.org/x/crypto/bcrypt"
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

// VerifySecret compares the given password with the stored password for the specified user and returns true if a match is found.
func (u *staticUserConfigurationProfileService) VerifySecret(
	username string,
	passwordBase64 string) (bool, error) {

	passwordBytes, decodeErr := base64.StdEncoding.DecodeString(passwordBase64)
	if decodeErr != nil {
		return false, fmt.Errorf("failed to decode password: %w", decodeErr)
	}

	profile, profileErr := u.GetProfile(username)
	if profileErr != nil {
		return false, fmt.Errorf("failed to get user profile: %w", profileErr)
	}

	err := bcrypt.CompareHashAndPassword([]byte(profile.Secret), passwordBytes)
	if err != nil {
		return false, fmt.Errorf("password invalid")
	}

	return true, nil
}

// VerifyPublicKey compares the given public key with stored keys for the specified user and returns true if a match is found.
func (u *staticUserConfigurationProfileService) VerifyPublicKey(
	username string,
	expectedKey ssh.PublicKey) (bool, error) {

	profile, err := u.GetProfile(username)
	if err != nil {
		return false, fmt.Errorf("failed to get user profile: %w", err)
	}

	expectedKeyBytes := expectedKey.Marshal()
	for _, formattedPublicKey := range profile.PublicKeys {
		key, _, _, _, sshParseKeyErr := ssh.ParseAuthorizedKey([]byte(formattedPublicKey))
		if sshParseKeyErr != nil {
			return false, fmt.Errorf("failed to parse public key: %w", sshParseKeyErr)
		}
		if bytes.Equal(key.Marshal(), expectedKeyBytes) {
			return true, nil
		}
	}

	return false, fmt.Errorf("no matching public key found")
}

// GetProfile retrieves the user profile for the given authenticated username from the configuration.
func (u *staticUserConfigurationProfileService) GetProfile(authenticatedUsername string) (configuration.UserProfile, error) {

	profile, found := u.settings.AuthServer.Users[authenticatedUsername]
	if !found {
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
