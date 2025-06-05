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

		registry.RegisterModule(modules.ApplicationConfigurationModule)
		registry.RegisterModule(modules.LoggingModule)
		registry.RegisterModule(modules.FiberModule)
		registry.RegisterModule(modules.RouteHandlersModule)
		registry.RegisterModule(modules.ApplicationServicesModule)

		return nil
	})

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
