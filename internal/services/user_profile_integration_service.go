package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

type userProfileIntegrationService struct {
	factory  func(string) (UserProfileBackend, error)
	settings *configuration.ApplicationConfiguration
}

func (s *userProfileIntegrationService) GetProfile(authenticatedUser string) (configuration.UserProfile, error) {
	for _, backendKey := range s.settings.AuthServer.Backends {
		backend, err := s.factory(backendKey)
		if err != nil || backend == nil {
			continue
		}

		profile, err := backend.GetProfile(authenticatedUser)
		if err != nil {
			if err.Error() == ErrorUserProfileNotFound {
				continue
			}
			return EmptyUserProfile, err
		}

		// Common normalization logic
		containerImageName := strings.TrimSpace(profile.Image)
		if containerImageName == "" {
			return EmptyUserProfile, errors.New(ErrorInvalidSessionContainerImage)
		}

		if len(profile.ShellCommand) == 0 {
			profile.ShellCommand = []string{"/bin/sh"}
		}

		return profile, nil
	}

	return EmptyUserProfile, errors.New(ErrorUserProfileNotFound)
}

func (s *userProfileIntegrationService) VerifyPublicKey(username string, key ssh.PublicKey) (bool, error) {
	profile, err := s.GetProfile(username)
	if err != nil {
		return false, err
	}

	expectedKeyBytes := key.Marshal()
	for _, formattedPublicKey := range profile.PublicKeys {
		parsedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(formattedPublicKey))
		if err != nil {
			return false, fmt.Errorf("failed to parse public key: %w", err)
		}
		if bytes.Equal(parsedKey.Marshal(), expectedKeyBytes) {
			return true, nil
		}
	}

	return false, fmt.Errorf("no matching public key found")
}

func (s *userProfileIntegrationService) VerifySecret(username string, passwordBase64 string) (bool, error) {
	passwordBytes, err := base64.StdEncoding.DecodeString(passwordBase64)
	if err != nil {
		return false, fmt.Errorf("failed to decode password: %w", err)
	}

	profile, err := s.GetProfile(username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(profile.Secret), passwordBytes)
	if err != nil {
		return false, fmt.Errorf("password invalid")
	}

	return true, nil
}

func (s *userProfileIntegrationService) RegisterUser(
	username string,
	profile configuration.UserProfile) error {
	for _, backendKey := range s.settings.AuthServer.Backends {
		backend, backendFactoryErr := s.factory(backendKey)
		if backendFactoryErr != nil || backend == nil {
			slog.Error("failed to create backend",
				slog.String("backendKey", backendKey),
				slog.Any("err", backendFactoryErr))
			continue
		}

		registerErr := backend.RegisterUser(username, profile)
		if registerErr != nil {
			if registerErr.Error() == ErrorReadOnlyBackend {
				continue
			}
			return fmt.Errorf("failed to register user with backend: %w", registerErr)
		}

		return nil
	}

	return errors.New("no writable backend found")
}

// NewUserProfileIntegrationService creates a new UserProfileService that integrates multiple backends.
func NewUserProfileIntegrationService(
	factory func(string) (UserProfileBackend, error),
	settings *configuration.ApplicationConfiguration) UserProfileService {
	return &userProfileIntegrationService{
		factory:  factory,
		settings: settings,
	}
}
