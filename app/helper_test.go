package app_test

import (
	"context"
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/cicd-lectures/vehicle-server/app"
	"github.com/cicd-lectures/vehicle-server/storage/vehiclestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

const (
	databaseName     = "vehicle-server"
	databaseUser     = "test"
	databasePassword = "test"
	listenAddress    = "127.0.0.1:0"
)

func setupEnvironment(t *testing.T) (*app.App, func()) {
	t.Helper()

	var (
		ctx    = context.Background()
		logger = zaptest.NewLogger(t)
	)

	// Start the database.
	database, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgis/postgis:16-3.4-alpine"),
		postgres.WithDatabase(databaseName),
		postgres.WithUsername(databaseUser),
		postgres.WithPassword(databasePassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
		withLogger(logger),
	)
	require.NoError(t, err)

	// Retrieve the database URL.
	dbURL, err := database.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	require.NoError(t, err)

	appExited := make(chan struct{})

	ctx, cancel := context.WithCancel(ctx)
	app, err := app.New(ctx, app.Config{DatabaseURL: dbURL, ListenAddress: listenAddress}, logger)
	require.NoError(t, err)

	// Start asynchronously the server.
	go func() {
		err := app.Run(ctx)
		if !assert.NoError(t, err) {
			// If we fail to run the app correctly, cancel everthing and prevent a dealock.
			cancel()
		}
		close(appExited)
	}()

	// Declares the teardown function,
	// done as an anonymous function to capture the necesary variables.
	tearDownEnvironment := func() {
		logger.Info("Cleaning up the test environment")

		cancel()

		<-appExited

		err := app.Close()
		assert.NoError(t, err)

		err = database.Stop(context.Background(), nil)
		assert.NoError(t, err)
	}

	// Poll the ready endpoint until the app is ready to accept requests.
	for {
		resp, err := http.Get("http://" + path.Join(app.ListenAddress(), "_", "ready"))
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		logger.Info("Request failed, retrying in one second", zap.Error(err))

		select {
		case <-ctx.Done():
			return app, tearDownEnvironment
		case <-time.After(100 * time.Millisecond):
			continue
		}
	}

	logger.Info("App is ready")

	return app, tearDownEnvironment
}

func withLogger(l *zap.Logger) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Logger = &loggerAdapter{logger: l.Sugar()}
	}
}

type loggerAdapter struct {
	logger *zap.SugaredLogger
}

func (l *loggerAdapter) Printf(template string, args ...any) {
	l.logger.Infof(template, args...)
}

func seedVehicles(t *testing.T, store vehiclestore.Store, vehicles ...vehiclestore.Vehicle) {
	t.Helper()
	ctx := context.Background()

	for _, v := range vehicles {
		_, err := store.Create(ctx, v)
		require.NoError(t, err)
	}
}
