package internal

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/matzefriedrich/containerssh-authserver/internal/resources"
	"github.com/rs/zerolog"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
)

type authServerApplication struct {
	app             *fiber.App
	config          *configuration.ApplicationConfiguration
	logger          *zerolog.Logger
	applicationInfo *configuration.ApplicationInfo
}

// Run starts the application server and listens on the configured port. Returns an error if the server fails to start.
func (a *authServerApplication) Run(_ context.Context) error {

	a.logger.Info().Msg("Starting authserver")
	a.logger.Info().Msgf("Version: %s", a.applicationInfo.VersionString())
	a.logger.Info().Msgf("Port: %d", a.config.Port)

	a.logger.Info().Msgf("Configuration path: %s", a.applicationInfo.ConfigurationPath)
	a.logger.Info().Msgf("Configuration type: %s", a.applicationInfo.ConfigurationType)

	a.printBanner()

	listenAddress := fmt.Sprintf(":%d", a.config.Port)
	return a.app.Listen(listenAddress)
}

func (a *authServerApplication) printBanner() {

	bannerFile, _ := resources.Resources.Open(resources.BannerTxt)
	defer bannerFile.Close()

	scanner := bufio.NewScanner(bannerFile)
	for scanner.Scan() {
		line := scanner.Text()
		a.logger.Info().Msg(line)
	}
}

var _ bootstrap.Application = (*authServerApplication)(nil)

// NewFiberApp initializes and configures a Fiber application with the provided route handlers, configuration, and logger.
func NewFiberApp(
	app *fiber.App,
	routeHandlers []handlers.RouteHandler,
	config *configuration.ApplicationConfiguration,
	applicationInfo *configuration.ApplicationInfo,
	logger *zerolog.Logger) bootstrap.Application {

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
	}))

	// Register RouteHandler services with the resolved Fiber instance.
	for _, routeHandler := range routeHandlers {
		routeHandler.Register(app)
	}

	return &authServerApplication{
		app:             app,
		config:          config,
		applicationInfo: applicationInfo,
		logger:          logger,
	}
}
