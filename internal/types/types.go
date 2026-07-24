package types

import (
	"context"

	"github.com/matzefriedrich/containerssh-authserver/internal/shims"
	"golang.org/x/crypto/ssh"
)

//go:generate parsley-cli generate mocks

const (
	ErrorUserProfileNotFound          = "user profile not found"
	ErrorInvalidSessionContainerImage = "invalid session container image"
	ErrorReadOnlyBackend              = "read-only backend"
)

var (
	EmptyUserProfile = UserProfile{}
	InvalidAppConfig shims.AppConfigShim
)

type UserProfile struct {
	Binds        []string `mapstructure:"binds"`
	Image        string   `mapstructure:"image"`
	Networks     []string `mapstructure:"networks"`
	PublicKeys   []string `mapstructure:"publicKeys"`
	ShellCommand []string `mapstructure:"shellCommand"`
	Secret       string   `mapstructure:"secret"`
}

type UserProfileService interface {
	GetProfile(ctx context.Context, authenticatedUser string) (UserProfile, error)
	VerifyPublicKey(ctx context.Context, username string, key ssh.PublicKey) (bool, error)
	VerifySecret(ctx context.Context, username string, passwordBase64 string) (bool, error)
	RegisterUser(ctx context.Context, username string, profile UserProfile) error
}

type UserProfileBackend interface {
	GetProfile(ctx context.Context, authenticatedUser string) (UserProfile, error)
	RegisterUser(ctx context.Context, username string, profile UserProfile) error
}

type ContainerAppConfigService interface {
	CreateApplicationConfigFor(ctx context.Context, authenticatedUser string) (shims.AppConfigShim, error)
}
