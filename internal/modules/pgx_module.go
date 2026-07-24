package modules

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"github.com/matzefriedrich/parsley/pkg/registration"
	parsleytypes "github.com/matzefriedrich/parsley/pkg/types"
)

func PgxModule(registry parsleytypes.ServiceRegistry) error {
	_ = registration.RegisterSingleton(registry, pgxPoolFactory)
	return nil
}

func pgxPoolFactory(ctx context.Context, settings *configuration.ApplicationConfiguration) (types.PgxPool, error) {
	postgres := settings.Databases.Postgres
	if postgres == nil {
		return nil, fmt.Errorf("no postgres database configuration found")
	}
	connectionString := postgres.ConnectionString()
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
