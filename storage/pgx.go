package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cicd-lectures/vehicle-server/storage/vehiclestore"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const createDBStatement = `
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE SCHEMA IF NOT EXISTS vehicle_server;
CREATE TABLE IF NOT EXISTS vehicle_server.vehicles (
	id SERIAL PRIMARY KEY,
	shortcode TEXT NOT NULL,
	battery SMALLINT,
	position GEOMETRY(POINT, 4326) not null
);
`

type PGXStore struct {
	conn *pgx.Conn
}

func NewPGXStore(ctx context.Context, databaseURL string, logger *zap.Logger) (*PGXStore, error) {
	var (
		conn *pgx.Conn
		err  error
	)

	err = retry(
		ctx,
		time.Second,
		10,
		func() error {
			logger.Info("Attempting to connect to the database")

			conn, err = pgx.Connect(ctx, databaseURL)
			if err != nil {
				return fmt.Errorf("could not connect to the database: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("could not ping the database: %w", err)
	}

	// Create the database. We should run migrations here,
	// but this is a toy project :-).
	if _, err := conn.Exec(ctx, createDBStatement); err != nil {
		return nil, fmt.Errorf("could not create the database: %w", err)
	}

	return &PGXStore{conn: conn}, nil
}

func (s *PGXStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.conn.Close(ctx); err != nil {
		return fmt.Errorf("could not close the connection to the database: %w", err)
	}

	return nil
}

func (s *PGXStore) Vehicle() vehiclestore.Store {
	return vehiclestore.NewPGXStore(s.conn)
}

func retry(ctx context.Context, retryInterval time.Duration, maxAttempts int, do func() error) error {
	var lastError error

	for i := 0; i < maxAttempts; i++ {
		err := do()
		if err == nil {
			return nil
		}
		if errors.Is(context.Canceled, err) {
			return err
		}

		lastError = err

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryInterval):

		}

	}

	return lastError
}
