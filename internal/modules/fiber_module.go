package modules

import (
	"github.com/gofiber/fiber/v3"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// FiberModule registers a Fiber application and configuration to the provided service registry.
func FiberModule(registry types.ServiceRegistry) error {

	const applicationName = "authserver"

	_ = registration.RegisterInstance(registry, fiber.Config{
		AppName:   applicationName,
		Immutable: true,
	})

	_ = registry.Register(newFiber, types.LifetimeSingleton)

	return nil
}

// newFiber Activator method for new Fiber instances.
func newFiber(config fiber.Config) *fiber.App {

	app := fiber.New(config)

	return app
}
