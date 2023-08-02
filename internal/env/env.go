package env

import (
	"context"

	"bookstore/internal/db"
)

type ServerEnv struct {
	db *db.DB
}

func New(ctx context.Context) {}
