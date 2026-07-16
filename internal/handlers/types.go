package handlers

import "github.com/go-chi/chi/v5"

// RouteHandler defines an interface for registering routes with a chi router.
type RouteHandler interface {
	Register(router chi.Router)
}
