package internal

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/matzefriedrich/containerssh-authserver/internal/resources"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers"

	"github.com/matzefriedrich/parsley/pkg/bootstrap"
)

type authServerApplication struct {
	app             *fiber.App
	config          *configuration.ApplicationConfiguration
	logger          *slog.Logger
	applicationInfo *configuration.ApplicationInfo
}

// Run starts the application server and listens on the configured port. Returns an error if the server fails to start.
func (a *authServerApplication) Run(ctx context.Context) error {

	a.printApplicationInfo()
	a.printBanner()

	listenAddress := fmt.Sprintf(":%d", a.config.Port)
	return a.app.Listen(listenAddress, fiber.ListenConfig{
		GracefulContext:       ctx,
		DisableStartupMessage: true,
	})
}

func (a *authServerApplication) printApplicationInfo() {
	a.logger.Info("Starting application", slog.String("releaseName", a.applicationInfo.ReleaseName))
	a.logger.Info("Version info", slog.String("version", a.applicationInfo.VersionString()))
	a.logger.Info("Network info", slog.Int("port", a.config.Port))
	a.logger.Info("Configuration info",
		slog.String("path", a.applicationInfo.ConfigurationPath),
		slog.String("type", a.applicationInfo.ConfigurationType))
}

func (a *authServerApplication) printBanner() {

	bannerFile, _ := resources.Resources.Open(resources.BannerTxt)
	defer func(bannerFile fs.File) {
		_ = bannerFile.Close()
	}(bannerFile)

	scanner := bufio.NewScanner(bannerFile)
	for scanner.Scan() {
		line := scanner.Text()
		a.logger.Info(line)
	}
}

var _ bootstrap.Application = (*authServerApplication)(nil)

// NewFiberApp initializes and configures a Fiber application with the provided route handlers, configuration, and logger.
func NewFiberApp(
	app *fiber.App,
	routeHandlers []handlers.RouteHandler,
	config *configuration.ApplicationConfiguration,
	applicationInfo *configuration.ApplicationInfo,
	logger *slog.Logger) bootstrap.Application {

	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		logger.Info("request",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("latency", latency),
		)
		return err
	})

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
