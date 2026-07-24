package apps

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/resources"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
)

type autoSetupApplication struct {
	config  *configuration.ApplicationConfiguration
	pool    types.PgxPool
	service types.UserProfileService
	logger  *slog.Logger
}

func (a *autoSetupApplication) Run(ctx context.Context) error {

	a.logger.Info("Starting auto-setup")

	migrateErr := a.migrateSchema(ctx)
	if migrateErr != nil {
		return fmt.Errorf("failed to migrate database schema: %w", migrateErr)
	}

	seedErr := a.seedUsers(ctx)
	if seedErr != nil {
		return seedErr
	}

	a.logger.Info("Auto-setup completed successfully")
	return nil
}

func (a *autoSetupApplication) seedUsers(ctx context.Context) error {
	a.logger.Info("Seeding users from static configuration")
	for username, profile := range a.config.AuthServer.Users {
		a.logger.Info("Registering user", slog.String("username", username))
		registerUsersErr := a.service.RegisterUser(ctx, username, profile)
		if registerUsersErr != nil {
			a.logger.Error("Failed to register user", slog.String("username", username), slog.Any("error", registerUsersErr))
			return fmt.Errorf("failed to register user: %w", registerUsersErr)
		}
	}
	return nil
}

func (a *autoSetupApplication) migrateSchema(ctx context.Context) error {
	schemaSql, _ := resources.Resources.ReadFile(resources.PostgresSchemaSql)
	a.logger.Info("Creating database schema")
	_, execErr := a.pool.Exec(ctx, string(schemaSql))
	if execErr != nil {
		return fmt.Errorf("failed to create database schema: %w", execErr)
	}
	return nil
}

func NewAutoSetupApp(
	config *configuration.ApplicationConfiguration,
	pool types.PgxPool,
	service types.UserProfileService,
	logger *slog.Logger) bootstrap.Application {
	return &autoSetupApplication{
		config:  config,
		pool:    pool,
		service: service,
		logger:  logger,
	}
}
