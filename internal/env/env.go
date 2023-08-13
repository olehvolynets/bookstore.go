package env

import (
	"context"

	"bookstore/internal/db"
)

type ServerEnv struct {
	db *db.DB // nolint:unused
}

func New(ctx context.Context) {}
