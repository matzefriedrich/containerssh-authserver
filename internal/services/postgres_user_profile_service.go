package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
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

	return entity.ToUserProfile(), nil
}

var _ UserProfileBackend = (*postgresUserProfileService)(nil)

// NewPostgresUserProfileService creates a new UserProfileBackend using a Postgres database.
func NewPostgresUserProfileService(pool postgresPool) UserProfileBackend {
	return &postgresUserProfileService{
		pool: pool,
	}
}
