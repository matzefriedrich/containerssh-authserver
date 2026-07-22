package services

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/shims"
	"golang.org/x/crypto/ssh"
)

//go:generate parsley-cli generate mocks

const (
	ErrorUserProfileNotFound          = "user profile not found"
	ErrorInvalidSessionContainerImage = "invalid session container image"
)

var (
	EmptyUserProfile = configuration.UserProfile{}
	InvalidAppConfig shims.AppConfigShim
)

type UserProfileService interface {
	GetProfile(authenticatedUser string) (configuration.UserProfile, error)
	VerifyPublicKey(username string, key ssh.PublicKey) (bool, error)
	VerifySecret(username string, passwordBase64 string) (bool, error)
}

type ContainerAppConfigService interface {
	CreateApplicationConfigFor(authenticatedUser string) (shims.AppConfigShim, error)
}
