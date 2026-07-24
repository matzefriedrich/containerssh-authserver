package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
)

type userProfileEntity struct {
	Binds        []string
	Image        string
	Networks     []string
	PublicKeys   []string
	ShellCommand []string
	Secret       string
}

func (e *userProfileEntity) ToUserProfile() types.UserProfile {
	return types.UserProfile{
		Binds:        e.Binds,
		Image:        e.Image,
		Networks:     e.Networks,
		PublicKeys:   e.PublicKeys,
		ShellCommand: e.ShellCommand,
		Secret:       e.Secret,
	}
}

type postgresUserProfileService struct {
	pool types.PgxPool
}

// GetProfile retrieves the user profile for the given authenticated username from the Postgres database.
func (p *postgresUserProfileService) GetProfile(
	ctx context.Context,
	authenticatedUsername string) (types.UserProfile, error) {

	query := `
		SELECT binds, image, networks, public_keys, shell_command, secret 
		FROM authserver.user_profiles 
		WHERE username = $1`

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
			return types.EmptyUserProfile, errors.New(types.ErrorUserProfileNotFound)
		}
		return types.EmptyUserProfile, fmt.Errorf("failed to query user profile: %w", err)
	}

	return entity.ToUserProfile(), nil
}

// RegisterUser adds or updates a user profile in the Postgres database.
func (p *postgresUserProfileService) RegisterUser(
	ctx context.Context,
	username string,
	profile types.UserProfile) error {

	query := `
		MERGE INTO authserver.user_profiles AS target
		USING (VALUES ($1::TEXT, $2::TEXT[], $3::TEXT, $4::TEXT[], $5::TEXT[], $6::TEXT[], $7::TEXT)) AS source(username, binds, image, networks, public_keys, shell_command, secret)
		ON (target.username = source.username)
		WHEN MATCHED THEN
			UPDATE SET
				binds = source.binds,
				image = source.image,
				networks = source.networks,
				public_keys = source.public_keys,
				shell_command = source.shell_command,
				secret = source.secret
		WHEN NOT MATCHED THEN
			INSERT (username, binds, image, networks, public_keys, shell_command, secret)
			VALUES (source.username, source.binds, source.image, source.networks, source.public_keys, source.shell_command, source.secret)`

	_, err := p.pool.Exec(ctx, query,
		username,
		profile.Binds,
		profile.Image,
		profile.Networks,
		profile.PublicKeys,
		profile.ShellCommand,
		profile.Secret,
	)

	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

var _ types.UserProfileBackend = (*postgresUserProfileService)(nil)

// NewPostgresUserProfileService creates a new UserProfileBackend using a Postgres database.
func NewPostgresUserProfileService(pool types.PgxPool) types.UserProfileBackend {
	return &postgresUserProfileService{
		pool: pool,
	}
}
