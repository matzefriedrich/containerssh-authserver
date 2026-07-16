package modules

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/logging"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// LoggingModule registers the required logging services.
func LoggingModule(registry types.ServiceRegistry) error {

	_ = registry.Register(logging.NewSlogLogger, types.LifetimeSingleton)

	return nil
}
