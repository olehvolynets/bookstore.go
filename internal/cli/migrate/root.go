package migrate

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"

	"bookstore/internal/db"
	"bookstore/pkg/migrate"
)

var Command = &cobra.Command{
	Use:   "migrate COMMAND",
	Short: "DB migrations management",
}

var migrEngine *migrate.Engine

func init() {
	ctx := context.Background()

	conf := db.Config{}
	if err := envconfig.Process(ctx, &conf); err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(ctx, conf.String())
	if err != nil {
		panic(err)
	}
	migrEngine = migrate.New(conn)

	Command.AddCommand(upCommand)
	Command.AddCommand(downCommand)
	Command.AddCommand(redoCommand)
	Command.AddCommand(createMigrationCommand)
	Command.AddCommand(statusCommand)
	Command.AddCommand(versionCommand)
}
