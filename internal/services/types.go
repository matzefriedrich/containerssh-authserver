package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
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
	EmptyUserProfile = configuration.UserProfile{}
	InvalidAppConfig shims.AppConfigShim
)

type UserProfileService interface {
	GetProfile(authenticatedUser string) (configuration.UserProfile, error)
	VerifyPublicKey(username string, key ssh.PublicKey) (bool, error)
	VerifySecret(username string, passwordBase64 string) (bool, error)
	RegisterUser(username string, profile configuration.UserProfile) error
}

type UserProfileBackend interface {
	GetProfile(authenticatedUser string) (configuration.UserProfile, error)
	RegisterUser(username string, profile configuration.UserProfile) error
}

type ContainerAppConfigService interface {
	CreateApplicationConfigFor(authenticatedUser string) (shims.AppConfigShim, error)
}

type PgxPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}
