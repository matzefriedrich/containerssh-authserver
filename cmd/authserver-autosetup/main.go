package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/modules"
	"github.com/matzefriedrich/containerssh-authserver/internal/resources"
	"github.com/matzefriedrich/containerssh-authserver/internal/services"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type autoSetupApplication struct {
	config  *configuration.ApplicationConfiguration
	pool    *pgxpool.Pool
	service services.UserProfileService
	logger  *slog.Logger
}

func (a *autoSetupApplication) Run(ctx context.Context) error {

	a.logger.Info("Starting auto-setup")

	// 1. Create schema
	schemaSql, err := resources.Resources.ReadFile(resources.PostgresSchemaSql)
	if err != nil {
		return err
	}

	a.logger.Info("Creating database schema")
	_, err = a.pool.Exec(ctx, string(schemaSql))
	if err != nil {
		return err
	}

	// 2. Seed users from static config
	a.logger.Info("Seeding users from static configuration")
	for username, profile := range a.config.AuthServer.Users {
		a.logger.Info("Registering user", slog.String("username", username))
		err = a.service.RegisterUser(username, profile)
		if err != nil {
			a.logger.Error("Failed to register user", slog.String("username", username), slog.Any("error", err))
			return err
		}
	}

	a.logger.Info("Auto-setup completed successfully")
	return nil
}

func NewAutoSetupApp(
	config *configuration.ApplicationConfiguration,
	pool *pgxpool.Pool,
	service services.UserProfileService,
	logger *slog.Logger) bootstrap.Application {
	return &autoSetupApplication{
		config:  config,
		pool:    pool,
		service: service,
		logger:  logger,
	}
}

func main() {

	ctx := context.Background()

	err := bootstrap.RunParsleyApplication(ctx, NewAutoSetupApp, func(registry types.ServiceRegistry) error {

		_ = registry.RegisterModule(modules.ApplicationConfigurationModule)
		_ = registry.RegisterModule(modules.LoggingModule)
		_ = registry.RegisterModule(modules.ApplicationServicesModule)

		return nil
	})

	if err != nil {
		slog.Error("Auto-setup failed", slog.Any("error", err))
		os.Exit(1)
	}

	os.Exit(0)
}
