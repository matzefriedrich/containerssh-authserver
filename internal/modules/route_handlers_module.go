package modules

import (
	"context"

	"github.com/matzefriedrich/containerssh-authserver/internal/handlers"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// RouteHandlersModule registers various route handlers in the provided service registry.
// It includes password, configuration, and public key handlers for route handling logic.
func RouteHandlersModule(registry types.ServiceRegistry) error {

	features.RegisterList[handlers.RouteHandler](context.Background(), registry)

	registration.RegisterTransient(registry, handlers.NewConfigHookHandler)
	registration.RegisterTransient(registry, handlers.NewPubKeyHookHandler)

	return nil
}
