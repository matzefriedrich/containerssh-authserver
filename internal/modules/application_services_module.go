package modules

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/services"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// ApplicationServicesModule registers application services with the provided service registry.
func ApplicationServicesModule(registry types.ServiceRegistry) error {

	registration.RegisterTransient(registry, services.NewStaticUserConfigurationProfileService)
	registration.RegisterTransient(registry, services.NewDockerAppConfigService)

	return nil
}
