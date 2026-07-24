package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

type userProfileIntegrationService struct {
	factory  func(string) (types.UserProfileBackend, error)
	settings *configuration.ApplicationConfiguration
}

func (s *userProfileIntegrationService) GetProfile(
	ctx context.Context,
	authenticatedUser string) (types.UserProfile, error) {
	for _, backendKey := range s.settings.AuthServer.Backends {
		backend, err := s.factory(backendKey)
		if err != nil || backend == nil {
			continue
		}

		profile, err := backend.GetProfile(ctx, authenticatedUser)
		if err != nil {
			if err.Error() == types.ErrorUserProfileNotFound {
				continue
			}
			return types.EmptyUserProfile, err
		}

		// Common normalization logic
		containerImageName := strings.TrimSpace(profile.Image)
		if containerImageName == "" {
			return types.EmptyUserProfile, errors.New(types.ErrorInvalidSessionContainerImage)
		}

		if len(profile.ShellCommand) == 0 {
			profile.ShellCommand = []string{"/bin/sh"}
		}

		return profile, nil
	}

	return types.EmptyUserProfile, errors.New(types.ErrorUserProfileNotFound)
}

func (s *userProfileIntegrationService) VerifyPublicKey(
	ctx context.Context,
	username string, key ssh.PublicKey) (bool, error) {
	profile, err := s.GetProfile(ctx, username)
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

func (s *userProfileIntegrationService) VerifySecret(
	ctx context.Context,
	username string,
	passwordBase64 string) (bool, error) {
	passwordBytes, err := base64.StdEncoding.DecodeString(passwordBase64)
	if err != nil {
		return false, fmt.Errorf("failed to decode password: %w", err)
	}

	profile, err := s.GetProfile(ctx, username)
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
	ctx context.Context,
	username string,
	profile types.UserProfile) error {
	for _, backendKey := range s.settings.AuthServer.Backends {
		backend, backendFactoryErr := s.factory(backendKey)
		if backendFactoryErr != nil || backend == nil {
			slog.Error("failed to create backend",
				slog.String("backendKey", backendKey),
				slog.Any("err", backendFactoryErr))
			continue
		}

		registerErr := backend.RegisterUser(ctx, username, profile)
		if registerErr != nil {
			if registerErr.Error() == types.ErrorReadOnlyBackend {
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
	factory func(string) (types.UserProfileBackend, error),
	settings *configuration.ApplicationConfiguration) types.UserProfileService {
	return &userProfileIntegrationService{
		factory:  factory,
		settings: settings,
	}
}
