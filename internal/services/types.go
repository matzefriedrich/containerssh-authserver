package services

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"

	"go.containerssh.io/containerssh/config"
	"golang.org/x/crypto/ssh"
)

//go:generate parsley-cli generate mocks

type UserProfileService interface {
	GetProfile(authenticatedUser string) (configuration.UserProfile, error)
	VerifyPublicKey(username string, key ssh.PublicKey) (bool, error)
}

var InvalidAppConfig config.AppConfig

type ContainerAppConfigService interface {
	CreateApplicationConfigFor(authenticatedUser string) (config.AppConfig, error)
}
