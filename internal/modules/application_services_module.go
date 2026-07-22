package modules

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/services"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// ApplicationServicesModule registers application services with the provided service registry.
func ApplicationServicesModule(registry types.ServiceRegistry) error {

	_ = registry.RegisterModule(PgxModule)
	_ = features.RegisterNamed[services.UserProfileBackend](registry,
		registration.NamedServiceRegistration("static", services.NewStaticUserConfigurationProfileService, types.LifetimeSingleton),
		registration.NamedServiceRegistration("postgres", services.NewPostgresUserProfileService, types.LifetimeScoped),
	)

	_ = registration.RegisterSingleton(registry, services.NewUserProfileIntegrationService)
	_ = registration.RegisterTransient(registry, services.NewDockerAppConfigService)

	return nil
}
