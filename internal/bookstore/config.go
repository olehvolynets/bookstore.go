package bookstore

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"

	"bookstore/internal/db"
)

type Config struct {
	Database db.Config

	Port uint   `env:"PORT, default=3000"`
	Env  string `env:"ENV, default=development"`
}

func NewConfig(ctx context.Context) (c Config, err error) {
	err = envconfig.Process(ctx, &c)
	if err != nil {
		return Config{}, fmt.Errorf("failed to setup server config: %w", err)
	}
	err = envconfig.Process(ctx, &c.Database)
	if err != nil {
		return Config{}, fmt.Errorf("failed to setup database config: %w", err)
	}

	return c, nil
}
