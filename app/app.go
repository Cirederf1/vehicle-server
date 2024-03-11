package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/cicd-lectures/vehicle-server/storage"
	"github.com/cicd-lectures/vehicle-server/vehicle"
	"go.uber.org/zap"
)

type App struct {
	listener net.Listener
	server   *http.Server
	store    *storage.PGXStore
	logger   *zap.Logger
}

type Config struct {
	DatabaseURL   string
	ListenAddress string
}

func New(ctx context.Context, cfg Config, logger *zap.Logger) (*App, error) {
	logger.Info(
		"Starting the vehicle-server",
		zap.String("database-url", cfg.DatabaseURL),
		zap.String("listen-address", cfg.ListenAddress),
	)

	// Initializing the storage layer.
	store, err := storage.NewPGXStore(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error(
			"Could not create the storage",
			zap.Error(err),
		)
		return nil, err
	}

	listener, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		logger.Error(
			"Could not listen",
			zap.String("listen-address", cfg.ListenAddress),
		)
		return nil, err
	}

	// Create up an http server and a router.
	var (
		router = http.NewServeMux()
		server = &http.Server{
			Addr:    cfg.ListenAddress,
			Handler: router,
		}
	)

	// Wire the routes.
	router.Handle("GET /vehicles", vehicle.NewListHandler(store, logger))
	router.Handle("POST /vehicles", vehicle.NewCreateHandler(store, logger))
	router.Handle("DELETE /vehicles/{id}", vehicle.NewDeleteHandler(store, logger))
	router.HandleFunc("GET /_/ready", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	return &App{
		listener: listener,
		server:   server,
		store:    store,
		logger:   logger,
	}, nil
}

func (a *App) ListenAddress() string {
	return a.listener.Addr().String()
}

func (a *App) Store() storage.Store {
	return a.store
}

func (a *App) Run(ctx context.Context) error {
	// Asynchronously watch for context cancelation.
	// And gracefully shutdown the server when it happens.
	// If the graceful shutdown happens, force it.
	go func() {
		<-ctx.Done()

		a.logger.Info("Shutting down the server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			a.logger.Error(
				"Could not gracefully shutdown the server, closing",
				zap.Error(err),
			)

			_ = a.server.Close()
		}
	}()

	a.logger.Info(
		"Listening for HTTP requests",
		zap.String("listen-address", a.listener.Addr().String()),
	)

	if err := a.server.Serve(a.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.Error(
			"Could not listen for HTTP requests",
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (a *App) Close() error {
	a.logger.Info("Server is stopping, see you next time!")

	if err := a.store.Close(); err != nil {
		a.logger.Error(
			"Could not close the storage",
			zap.Error(err),
		)
		return err
	}

	return nil
}
