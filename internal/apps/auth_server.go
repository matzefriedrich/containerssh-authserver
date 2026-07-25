package apps

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/matzefriedrich/containerssh-authserver/internal/resources"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
	"github.com/matzefriedrich/containerssh-authserver/internal/handlers"

	"github.com/matzefriedrich/parsley/pkg/bootstrap"
)

type authServerApplication struct {
	router          *chi.Mux
	config          *configuration.ApplicationConfiguration
	logger          *slog.Logger
	applicationInfo *configuration.ApplicationInfo
}

// Run starts the application server and listens on the configured port. Returns an error if the server fails to start.
func (a *authServerApplication) Run(ctx context.Context) error {

	a.printApplicationInfo()
	a.printBanner()

	listenAddress := fmt.Sprintf(":%d", a.config.Port)
	server := &http.Server{
		Addr:    listenAddress,
		Handler: a.router,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
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

// NewChiApp initializes and configures a chi application with the provided route handlers, configuration, and logger.
func NewChiApp(
	router *chi.Mux,
	routeHandlers []handlers.RouteHandler,
	config *configuration.ApplicationConfiguration,
	applicationInfo *configuration.ApplicationInfo,
	logger *slog.Logger) bootstrap.Application {

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			latency := time.Since(start)

			logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Duration("latency", latency),
			)
		})
	})

	// Register RouteHandler services with the resolved chi instance.
	for _, routeHandler := range routeHandlers {
		routeHandler.Register(router)
	}

	return &authServerApplication{
		router:          router,
		config:          config,
		applicationInfo: applicationInfo,
		logger:          logger,
	}
}
