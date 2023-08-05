package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"bookstore/internal/logging"
)

type DB struct {
	Pool *pgxpool.Pool
	Log  *zerolog.Logger
}

func New(ctx context.Context, cfg *Config) (*DB, error) {
	globalLogger := logging.FromContext(ctx)
	logger := logging.WithLogSource(globalLogger, "db")

	pgxConfig, err := pgxpool.ParseConfig(cfg.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pgxConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// Ping the connection to see if it is still valid. Ping returns an error if
		// it fails.
		return conn.Ping(ctx) == nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	logger.Info().Msg("creating connection pool")

	return &DB{pool, &logger}, nil
}

func (db *DB) Close(_ context.Context) {
	db.Log.Info().Msg("closing connection pool")
	db.Pool.Close()
}
