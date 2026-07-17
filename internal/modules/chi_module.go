package modules

import (
	"github.com/go-chi/chi/v5"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// ChiModule registers a chi router to the provided service registry.
func ChiModule(registry types.ServiceRegistry) error {

	_ = registry.Register(chi.NewRouter, types.LifetimeSingleton)

	return nil
}
