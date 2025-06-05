package modules

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/logging"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// LoggingModule registers the required logging services.
func LoggingModule(registry types.ServiceRegistry) error {

	registry.Register(logging.NewZeroLogLogger, types.LifetimeSingleton)

	return nil
}
