package main

import (
	"context"
	"os"

	"github.com/matzefriedrich/containerssh-authserver/internal"
	"github.com/matzefriedrich/containerssh-authserver/internal/modules"

	"github.com/matzefriedrich/parsley/pkg/bootstrap"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func main() {

	ctx := context.Background()

	err := bootstrap.RunParsleyApplication(ctx, internal.NewFiberApp, func(registry types.ServiceRegistry) error {

		_ = registry.RegisterModule(modules.ApplicationConfigurationModule)
		_ = registry.RegisterModule(modules.LoggingModule)
		_ = registry.RegisterModule(modules.FiberModule)
		_ = registry.RegisterModule(modules.RouteHandlersModule)
		_ = registry.RegisterModule(modules.ApplicationServicesModule)

		return nil
	})

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
