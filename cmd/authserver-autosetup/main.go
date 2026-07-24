package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/matzefriedrich/containerssh-authserver/internal/apps"
	"github.com/matzefriedrich/containerssh-authserver/internal/modules"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
	parsleytypes "github.com/matzefriedrich/parsley/pkg/types"
)

func main() {

	ctx := context.Background()

	err := bootstrap.RunParsleyApplication(ctx, apps.NewAutoSetupApp, func(registry parsleytypes.ServiceRegistry) error {

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
