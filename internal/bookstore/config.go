package bookstore

import (
	"context"

	"github.com/sethvargo/go-envconfig"

	"bookstore/internal/db"
)

type Config struct {
	Database *db.Config

	Port uint   `env:"PORT, default=3000"`
	Env  string `env:"ENV, default=development"`
}

func NewConfig(ctx context.Context, dbConfig *db.Config) *Config {
	c := &Config{
		Database: dbConfig,
	}

	envconfig.Process(ctx, c)

	return c
}
