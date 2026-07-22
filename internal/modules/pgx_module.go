package modules

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func PgxModule(registry types.ServiceRegistry) error {
	_ = registration.RegisterSingleton(registry, pgxPoolFactory)
	return nil
}

func pgxPoolFactory(ctx context.Context, settings *configuration.ApplicationConfiguration) (*pgxpool.Pool, error) {
	postgres := settings.Databases.Postgres
	if postgres == nil {
		return nil, fmt.Errorf("no postgres database configuration found")
	}
	connectionString := postgres.ConnectionString()
	return pgxpool.New(ctx, connectionString)
}
