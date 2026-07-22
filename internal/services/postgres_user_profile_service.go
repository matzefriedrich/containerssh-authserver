package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

type userProfileEntity struct {
	Binds        []string
	Image        string
	Networks     []string
	PublicKeys   []string
	ShellCommand []string
	Secret       string
}

func (e *userProfileEntity) ToUserProfile() configuration.UserProfile {
	return configuration.UserProfile{
		Binds:        e.Binds,
		Image:        e.Image,
		Networks:     e.Networks,
		PublicKeys:   e.PublicKeys,
		ShellCommand: e.ShellCommand,
		Secret:       e.Secret,
	}
}

type postgresPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type postgresUserProfileService struct {
	pool postgresPool
}

// Key returns the unique identifier for this user profile service implementation.
func (p *postgresUserProfileService) Key() string {
	return "postgres"
}

// VerifySecret compares the given password with the stored password for the specified user and returns true if a match is found.
func (p *postgresUserProfileService) VerifySecret(
	username string,
	passwordBase64 string) (bool, error) {

	passwordBytes, decodeErr := base64.StdEncoding.DecodeString(passwordBase64)
	if decodeErr != nil {
		return false, fmt.Errorf("failed to decode password: %w", decodeErr)
	}

	profile, profileErr := p.GetProfile(username)
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
func (p *postgresUserProfileService) VerifyPublicKey(
	username string,
	expectedKey ssh.PublicKey) (bool, error) {

	profile, err := p.GetProfile(username)
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

// GetProfile retrieves the user profile for the given authenticated username from the Postgres database.
func (p *postgresUserProfileService) GetProfile(authenticatedUsername string) (configuration.UserProfile, error) {

	ctx := context.Background()
	query := "SELECT binds, image, networks, public_keys, shell_command, secret FROM authserver.user_profiles WHERE username = $1"

	var entity userProfileEntity
	err := p.pool.QueryRow(ctx, query, authenticatedUsername).Scan(
		&entity.Binds,
		&entity.Image,
		&entity.Networks,
		&entity.PublicKeys,
		&entity.ShellCommand,
		&entity.Secret,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return EmptyUserProfile, errors.New(ErrorUserProfileNotFound)
		}
		return EmptyUserProfile, fmt.Errorf("failed to query user profile: %w", err)
	}

	profile := entity.ToUserProfile()

	containerImageName := strings.TrimSpace(profile.Image)
	if containerImageName == "" {
		return EmptyUserProfile, errors.New(ErrorInvalidSessionContainerImage)
	}

	if len(profile.ShellCommand) == 0 {
		profile.ShellCommand = []string{"/bin/sh"}
	}

	return profile, nil
}

var _ UserProfileService = (*postgresUserProfileService)(nil)

// NewPostgresUserProfileService creates a new UserProfileService using a Postgres database.
func NewPostgresUserProfileService(pool postgresPool) UserProfileService {
	return &postgresUserProfileService{
		pool: pool,
	}
}
