package services

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/types/shims"
	"golang.org/x/crypto/ssh"
)

//go:generate parsley-cli generate mocks

type UserProfileService interface {
	GetProfile(authenticatedUser string) (configuration.UserProfile, error)
	VerifyPublicKey(username string, key ssh.PublicKey) (bool, error)
}

var InvalidAppConfig shims.AppConfigShim

type ContainerAppConfigService interface {
	CreateApplicationConfigFor(authenticatedUser string) (shims.AppConfigShim, error)
}
