package handlers

import "github.com/gofiber/fiber/v2"

// RouteHandler defines an interface for registering routes with a Fiber application.
type RouteHandler interface {
	Register(app *fiber.App)
}
