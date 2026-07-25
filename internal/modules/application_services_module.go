package modules

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/services"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	parsleytypes "github.com/matzefriedrich/parsley/pkg/types"
)

// ApplicationServicesModule registers application services with the provided service registry.
func ApplicationServicesModule(registry parsleytypes.ServiceRegistry) error {

	_ = registry.RegisterModule(PgxModule)
	_ = features.RegisterNamed[types.UserProfileBackend](registry,
		registration.NamedServiceRegistration("static", services.NewStaticUserConfigurationProfileService, parsleytypes.LifetimeSingleton),
		registration.NamedServiceRegistration("postgres", services.NewPostgresUserProfileService, parsleytypes.LifetimeScoped),
	)

	_ = registration.RegisterSingleton(registry, services.NewUserProfileIntegrationService)
	_ = registration.RegisterTransient(registry, services.NewDockerAppConfigService)

	return nil
}
